package sniffer

import (
	"context"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/RyanCarrier/dijkstra"

	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "my.domain/Len/api/v1"
	"my.domain/Len/pkg/log"
)

type Sniffer struct {
	// CRD Info
	cache  cache.Cache
	client client.Client

	nodeName string

	// Edge Network Graph File
	fileName string

	// Edge Node Name
	nodeNameList []string

	// Latency Info
	latencyList v1.LatencyList

	updateInterval int64
}

func (s *Sniffer) getNodes() {
	newNodeNameList := make([]string, 0)
	got, err := os.ReadFile(s.fileName)
	if err != nil {
		log.ErrPrint(err)
	}
	input := strings.TrimSpace(string(got))
	for _, line := range strings.Split(input, "\n") {
		f := strings.Fields(strings.TrimSpace(line))
		if len(f) == 0 || len(f) == 1 {
			continue
		}
		newNodeNameList = append(newNodeNameList, f[0])
	}
	if len(s.nodeNameList) == 0 || !reflect.DeepEqual(s.nodeNameList, newNodeNameList) {
		s.nodeNameList = newNodeNameList
	}
	//s.nodeNameList = newNodeNameList
}

func (s *Sniffer) updateLatency() {
	newLatencyList := make(v1.LatencyList, 0)
	graph, err := dijkstra.Import(s.fileName)
	s.nodeName = os.Getenv("NODE_NAME")
	if err != nil {
		log.ErrPrint(err)
	}
	s.getNodes()
	i, _ := graph.GetMapping(s.nodeName)
	for _, n := range s.nodeNameList {
		if n == s.nodeName {
			continue
		}
		j, _ := graph.GetMapping(n)
		best, err1 := graph.Shortest(i, j)
		if err1 != nil {
			log.ErrPrint(err)
		}
		newLatencyList = append(newLatencyList, v1.Latency{
			NodeName: n,
			Latency:  best.Distance,
		})
	}
	if len(s.latencyList) == 0 || !reflect.DeepEqual(s.latencyList, newLatencyList) {
		s.latencyList = newLatencyList
	}
	//s.latencyList = newLatencyList
}

func (s *Sniffer) createLen() error {
	len := v1.Len{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.nodeName,
		},
		Spec: v1.LenSpec{
			UpdateInterval: s.updateInterval,
		},
	}
	err := s.client.Create(context.TODO(), &len)
	if err != nil && !apierror.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (s *Sniffer) Process() {
	interval := time.Duration(s.updateInterval) * time.Millisecond
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		s.updateLatency()

		currentLen := v1.Len{}

		key := types.NamespacedName{
			Name: s.nodeName,
		}

		err := s.client.Get(context.TODO(), key, &currentLen)
		if err != nil {
			log.ErrPrint(err)
			continue
		}
		if s.NeedUpdate(currentLen.Status) {
			updateLen := currentLen.DeepCopy()
			updateLen.Status = v1.LenStatus{
				NodeName:    s.nodeName,
				LatencyList: s.latencyList,
				UpdateTime: &metav1.Time{
					Time: time.Now(),
				},
			}

			if err := s.client.Update(context.TODO(), updateLen); err != nil {
				log.ErrPrint(err)
			}
		}
	}
}

func NewSniffer(filePath string, interval int64, client client.Client, cache cache.Cache) *Sniffer {
	return &Sniffer{
		latencyList:    make(v1.LatencyList, 0),
		fileName:       filePath,
		nodeNameList:   []string{},
		updateInterval: interval,
		client:         client,
		cache:          cache,
	}
}

func StartSniffer(s *Sniffer) {
	// Init CRD & Set Config
	s.nodeName = os.Getenv("NODE_NAME")
	if err := s.createLen(); err != nil {
		panic(err)
	}
	s.Process()
}

func (s *Sniffer) NeedUpdate(status v1.LenStatus) bool {
	if status.UpdateTime == nil {
		log.Print("LatencyList is Null, needs update.")
		return true
	}

	if status.NodeName != s.nodeName {
		log.Print("Edge Node Name changed, needs update.")
		return true
	}

	if !reflect.DeepEqual(status.LatencyList, s.latencyList) {
		log.Print("Latency List changed, needs update.")
		return true
	}
	return false
}
