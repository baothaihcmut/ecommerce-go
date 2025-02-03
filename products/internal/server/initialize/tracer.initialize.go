package initialize

import (
	appCfg "github.com/baothaihcmut/Ecommerce-Go/products/internal/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func InitializeTracer(cfg *appCfg.Config) (opentracing.Tracer, error) {
	jaegerCfg := &config.Configuration{
		ServiceName: "Product service",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: cfg.Jaeger.Address,
		},
	}
	tracer, closer, err := jaegerCfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	return tracer, nil
}
