package dynatrace

import (
	"fmt"
	"strings"
	"time"
)

const timeFormat = "2006-01-02T15:04:05Z"

type DTQuery struct {
	fragments  []string
	finalQuery string
}

func (q *DTQuery) InitLogs(hours int) *DTQuery {
	q.fragments = []string{}

	q.fragments = append(q.fragments, fmt.Sprintf("fetch logs, from:now()-%dh \n| filter matchesValue(event.type, \"LOG\") and ", hours))

	return q
}

func (q *DTQuery) InitLogsWithTimeRange(from time.Time, to time.Time) *DTQuery {
	q.fragments = []string{}

	fromStr := from.Format(timeFormat)
	toStr := to.Format(timeFormat)

	q.fragments = append(q.fragments, fmt.Sprintf("fetch logs, from:\"%s\", to:\"%s\" \n| filter matchesValue(event.type, \"LOG\") and ", fromStr, toStr))

	return q
}

func (q *DTQuery) InitEvents(hours int) *DTQuery {
	q.fragments = []string{}

	q.fragments = append(q.fragments, fmt.Sprintf("fetch events, from:now()-%dh \n| filter ", hours))

	return q
}

func (q *DTQuery) Cluster(mgmtClusterName string) *DTQuery {
	q.fragments = append(q.fragments, fmt.Sprintf("matchesPhrase(dt.kubernetes.cluster.name, \"%s\")", mgmtClusterName))

	return q
}

func (q *DTQuery) Namespaces(namespaceList []string) *DTQuery {
	var nsQuery string
	finalQuery := ""
	nsQuery = " and ("

	for i, ns := range namespaceList {
		nsQuery += fmt.Sprintf("matchesValue(k8s.namespace.name, \"%s\")", ns)
		if i < len(namespaceList)-1 {
			nsQuery += " or "
		}
	}
	nsQuery += ")"
	finalQuery += nsQuery

	q.fragments = append(q.fragments, finalQuery)

	return q
}

func (q *DTQuery) Nodes(nodeList []string) *DTQuery {
	var nodeQuery string

	nodeQuery = " and ("
	for i, node := range nodeList {
		nodeQuery += fmt.Sprintf("matchesValue(k8s.node.name, \"%s\")", node)
		if i < len(nodeList)-1 {
			nodeQuery += " or "
		}
	}
	nodeQuery += ")"
	q.fragments = append(q.fragments, nodeQuery)

	return q
}

func (q *DTQuery) Pods(podList []string) *DTQuery {
	var podQuery string

	podQuery = " and ("
	for i, pod := range podList {
		podQuery += fmt.Sprintf("matchesValue(k8s.pod.name, \"%s\")", pod)
		if i < len(podList)-1 {
			podQuery += " or "
		}
	}
	podQuery += ")"
	q.fragments = append(q.fragments, podQuery)

	return q
}

func (q *DTQuery) Containers(containerList []string) *DTQuery {
	var containerQuery string

	containerQuery = " and ("
	for i, container := range containerList {
		containerQuery += fmt.Sprintf("matchesValue(k8s.container.name, \"%s\")", container)
		if i < len(containerList)-1 {
			containerQuery += " or "
		}
	}
	containerQuery += ")"
	q.fragments = append(q.fragments, containerQuery)

	return q
}

func (q *DTQuery) Status(statusList []string) *DTQuery {
	var statusQuery string

	statusQuery = " and ("
	for i, status := range statusList {
		statusQuery += fmt.Sprintf("matchesValue(status, \"%s\")", status)
		if i < len(statusList)-1 {
			statusQuery += " or "
		}
	}
	statusQuery += ")"
	q.fragments = append(q.fragments, statusQuery)

	return q
}

func (q *DTQuery) ContainsPhrase(phrase string) *DTQuery {
	q.fragments = append(q.fragments, " and contains(content,\""+phrase+"\", caseSensitive:false)")

	return q
}

func (q *DTQuery) Sort(order string) (query *DTQuery, error error) {
	validOrders := []string{
		"asc",
		"desc",
	}

	for _, or := range validOrders {
		if or == order {
			q.fragments = append(q.fragments, fmt.Sprintf("\n| sort timestamp %s", order))
			return q, nil
		}
	}

	return q, fmt.Errorf("no valid sorting order specified. valid order are %s. given %v", strings.Join(validOrders, ", "), order)
}

func (q *DTQuery) Deployments(workloads []string) *DTQuery {
	var deploymentQuery string

	deploymentQuery = " and ("
	for i, deploy := range workloads {
		deploymentQuery += fmt.Sprintf("matchesValue(dt.kubernetes.workload.name, \"%s\")", deploy)
		if i < len(workloads)-1 {
			deploymentQuery += " or "
		}
	}
	deploymentQuery += ")"
	q.fragments = append(q.fragments, deploymentQuery)

	return q
}

func (q *DTQuery) Limit(limit int) *DTQuery {
	q.fragments = append(q.fragments, "\n| limit "+fmt.Sprint(limit))

	return q
}

func (q *DTQuery) Build() string {
	q.finalQuery = strings.Join(q.fragments[:], "")

	return q.finalQuery
}

type ManagedQuery struct {
	filters []string
	query   string
}

func (mq *ManagedQuery) Cluster(clusterName string) *ManagedQuery {
	mq.filters = append(mq.filters, fmt.Sprintf("dt.kubernetes.cluster.name=\"%s\"", clusterName))
	return mq
}

func (mq *ManagedQuery) Namespaces(namespaces []string) *ManagedQuery {
	if len(namespaces) == 1 {
		mq.filters = append(mq.filters, fmt.Sprintf("k8s.namespace.name=\"%s\"", namespaces[0]))
	} else if len(namespaces) > 1 {
		parts := make([]string, len(namespaces))
		for i, ns := range namespaces {
			parts[i] = fmt.Sprintf("k8s.namespace.name=\"%s\"", ns)
		}
		mq.filters = append(mq.filters, "("+strings.Join(parts, " OR ")+")")
	}
	return mq
}

func (mq *ManagedQuery) Pods(pods []string) *ManagedQuery {
	if len(pods) == 1 {
		mq.filters = append(mq.filters, fmt.Sprintf("k8s.pod.name=\"%s\"", pods[0]))
	} else if len(pods) > 1 {
		parts := make([]string, len(pods))
		for i, p := range pods {
			parts[i] = fmt.Sprintf("k8s.pod.name=\"%s\"", p)
		}
		mq.filters = append(mq.filters, "("+strings.Join(parts, " OR ")+")")
	}
	return mq
}

func (mq *ManagedQuery) Nodes(nodes []string) *ManagedQuery {
	if len(nodes) == 1 {
		mq.filters = append(mq.filters, fmt.Sprintf("k8s.node.name=\"%s\"", nodes[0]))
	} else if len(nodes) > 1 {
		parts := make([]string, len(nodes))
		for i, n := range nodes {
			parts[i] = fmt.Sprintf("k8s.node.name=\"%s\"", n)
		}
		mq.filters = append(mq.filters, "("+strings.Join(parts, " OR ")+")")
	}
	return mq
}

func (mq *ManagedQuery) Containers(containers []string) *ManagedQuery {
	if len(containers) == 1 {
		mq.filters = append(mq.filters, fmt.Sprintf("k8s.container.name=\"%s\"", containers[0]))
	} else if len(containers) > 1 {
		parts := make([]string, len(containers))
		for i, c := range containers {
			parts[i] = fmt.Sprintf("k8s.container.name=\"%s\"", c)
		}
		mq.filters = append(mq.filters, "("+strings.Join(parts, " OR ")+")")
	}
	return mq
}

func (mq *ManagedQuery) Status(statuses []string) *ManagedQuery {
	if len(statuses) == 1 {
		mq.filters = append(mq.filters, fmt.Sprintf("status=\"%s\"", statuses[0]))
	} else if len(statuses) > 1 {
		parts := make([]string, len(statuses))
		for i, s := range statuses {
			parts[i] = fmt.Sprintf("status=\"%s\"", s)
		}
		mq.filters = append(mq.filters, "("+strings.Join(parts, " OR ")+")")
	}
	return mq
}

func (mq *ManagedQuery) ContainsPhrase(phrase string) *ManagedQuery {
	mq.filters = append(mq.filters, fmt.Sprintf("content=\"%s\"", phrase))
	return mq
}

func (mq *ManagedQuery) Deployments(workloads []string) *ManagedQuery {
	if len(workloads) == 1 {
		mq.filters = append(mq.filters, fmt.Sprintf("dt.kubernetes.workload.name=\"%s\"", workloads[0]))
	} else if len(workloads) > 1 {
		parts := make([]string, len(workloads))
		for i, w := range workloads {
			parts[i] = fmt.Sprintf("dt.kubernetes.workload.name=\"%s\"", w)
		}
		mq.filters = append(mq.filters, "("+strings.Join(parts, " OR ")+")")
	}
	return mq
}

func (mq *ManagedQuery) Build() string {
	mq.query = strings.Join(mq.filters, " AND ")
	return mq.query
}
