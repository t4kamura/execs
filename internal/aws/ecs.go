package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/aws/session-manager-plugin/src/log"
	"github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession"
	"github.com/google/uuid"
)

const (
	MAX_RESULTS_CLUSTER = 100
	MAX_RESULTS_SERVICE = 100
	MAX_RESULTS_TASK    = 100
)

type ECS struct {
	clnt *ecs.Client
}

func New(cfg *aws.Config) *ECS {
	return &ECS{
		clnt: ecs.NewFromConfig(*cfg),
	}
}

func (e *ECS) ListClusterNames(ctx context.Context) ([]string, error) {
	resp, err := e.clnt.ListClusters(ctx, &ecs.ListClustersInput{
		MaxResults: aws.Int32(MAX_RESULTS_CLUSTER),
	})
	if err != nil {
		return nil, err
	}
	names, err := extractNameFromARNs(resp.ClusterArns)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (e *ECS) ListServiceNames(ctx context.Context, cluster string) ([]string, error) {
	resp, err := e.clnt.ListServices(ctx, &ecs.ListServicesInput{
		Cluster:    aws.String(cluster),
		MaxResults: aws.Int32(MAX_RESULTS_SERVICE),
	})
	if err != nil {
		return nil, err
	}
	names, err := extractNameFromARNs(resp.ServiceArns)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (e *ECS) ListTaskNames(ctx context.Context, cluster string, service string) ([]string, error) {
	resp, err := e.clnt.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:       aws.String(cluster),
		ServiceName:   aws.String(service),
		MaxResults:    aws.Int32(MAX_RESULTS_TASK),
		DesiredStatus: "RUNNING",
	})
	if err != nil {
		return nil, err
	}
	names, err := extractNameFromARNs(resp.TaskArns)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (e *ECS) ListContainerNames(ctx context.Context, cluster string, task string) ([]string, error) {
	var names []string
	resp, err := e.clnt.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Tasks:   []string{task},
		Cluster: aws.String(cluster),
	})
	if err != nil {
		return names, err
	}

	if len(resp.Tasks) == 0 {
		return names, nil
	}

	for _, container := range resp.Tasks[0].Containers {
		names = append(names, *container.Name)
	}

	return names, nil
}

func (e *ECS) StartSession(ctx context.Context, cluster string, task string, container string, shell string, region string) error {
	resp, err := e.clnt.ExecuteCommand(ctx, &ecs.ExecuteCommandInput{
		Command:     aws.String(shell),
		Interactive: true,
		Task:        aws.String(task),
		Cluster:     aws.String(cluster),
		Container:   aws.String(container),
	})
	if err != nil {
		return err
	}

	endpoint, err := ssm.NewDefaultEndpointResolver().ResolveEndpoint(region, ssm.EndpointResolverOptions{})
	if err != nil {
		return err
	}

	sess := session.Session{
		SessionId:   aws.ToString(resp.Session.SessionId),
		StreamUrl:   aws.ToString(resp.Session.StreamUrl),
		TokenValue:  aws.ToString(resp.Session.TokenValue),
		Endpoint:    endpoint.URL,
		ClientId:    uuid.NewString(),
		DataChannel: &datachannel.DataChannel{},
	}

	return sess.Execute(log.Logger(false, sess.ClientId))
}
