package surveymarkup

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"project-survey-proceeder/internal/surveymarkup/contracts"
	"project-survey-proceeder/internal/surveymarkup/model"
	"project-survey-proceeder/internal/utils"
	"time"

	"google.golang.org/grpc"
)

type Service struct {
	srvGeneratorAddr string

	conn   *grpc.ClientConn
	client model.SurveyMarkupGeneratorClient
}

func NewService(surveyGeneratorAddress string) contracts.ISurveyMarkupService {
	return &Service{srvGeneratorAddr: surveyGeneratorAddress}
}

func (s *Service) Init() error {
	var err error
	s.conn, err = grpc.Dial(s.srvGeneratorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to SurveyMarkupGenerator: %v", err)
	}

	s.client = model.NewSurveyMarkupGeneratorClient(s.conn)
	return nil
}

func (s *Service) GetMarkup(unitId int, surveyIds []int, language string) (string, error) {
	conn, err := grpc.Dial(s.srvGeneratorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := model.NewSurveyMarkupGeneratorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GenerateMarkup(ctx, &model.GenerateMarkupRequest{
		UnitId:    int32(unitId),
		SurveyIds: utils.ConvertInts[int32](surveyIds),
		Language:  language,
	})

	if err != nil {
		return "", fmt.Errorf("could not get markup: %v", err)
	}

	return r.GetMarkup(), nil
}

func (s *Service) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}
