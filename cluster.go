package swarmtool

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"log"
)

type (
	// Node represents a swarm node.
	Node struct {
		Reachability string
	}
	// DockerClient is the API client that performs all operations
	// against a docker server.
	DockerClient interface {
		ManagerList(ctx context.Context) ([]*Node, error)
	}
	// Cluster represents a cluster of swarm nodes.
	Cluster struct {
		Client DockerClient
	}
	// ClusterClient implements DockerClient for cluster operations.
	ClusterClient struct {
		*client.Client
	}
)

// ManagerList returns list of manager nodes.
func (c *ClusterClient) ManagerList(ctx context.Context) ([]*Node, error) {
	fs := filters.NewArgs(filters.KeyValuePair{
		Key:   "role",
		Value: string(swarm.NodeRoleManager),
	})
	dockerNodes, err := c.NodeList(ctx, types.NodeListOptions{Filters: fs})
	if err != nil {
		return nil, err
	}

	var nodes []*Node
	for _, node := range dockerNodes {
		nodes = append(nodes, &Node{Reachability: string(node.ManagerStatus.Reachability)})
	}

	return nodes, nil
}

// IsActive checks if a node is reachable or not.
func (n *Node) IsActive() bool {
	if n.Reachability == "reachable" {
		return true
	}
	return false
}

// IsSafeToShutdown checks if cluster stays functional after
// shutting down a manager node.
func (c *Cluster) IsSafeToShutdown() bool {
	ctx := context.Background()
	allManagers, err := c.Client.ManagerList(ctx)
	if err != nil {
		log.Printf("failed to get swarm managers with error %s", err)
		return false
	}

	var availableManagers []*Node
	for _, manager := range allManagers {
		if manager.IsActive() {
			availableManagers = append(availableManagers, manager)
		}
	}

	return c.staysMajority(len(allManagers), len(availableManagers))
}

func (c *Cluster) staysMajority(all, available int) bool {
	if available-1 >= all {
		log.Printf("safe to shutdown node. managers: %d availabel: %d", all, available)
		return true
	}
	log.Printf("not safe to shutdown node. managers: %d availabel: %d", all, available)
	return false
}
