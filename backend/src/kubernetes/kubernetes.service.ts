import { Injectable, Logger } from '@nestjs/common';
import * as k8s from '@kubernetes/client-node';

@Injectable()
export class KubernetesService {
  private readonly logger = new Logger(KubernetesService.name);
  private batchV1Api: k8s.BatchV1Api;

  constructor() {
    const kc = new k8s.KubeConfig();
    try {
      kc.loadFromDefault();
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : String(err);
      this.logger.error(
        'Erro ao carregar configuração do Kubernetes',
        errorMessage,
      );
    }
    this.batchV1Api = kc.makeApiClient(k8s.BatchV1Api);
  }

  async createWorkerJob(operationId: string): Promise<void> {
    const job: k8s.V1Job = {
      apiVersion: 'batch/v1',
      kind: 'Job',
      metadata: {
        name: `worker-job-${operationId}`,
      },
      spec: {
        ttlSecondsAfterFinished: 30,
        template: {
          metadata: {
            labels: {
              app: 'worker',
            },
          },
          spec: {
            restartPolicy: 'Never',
            containers: [
              {
                name: 'worker',
                image: 'worker_v1',
                imagePullPolicy: 'Never',
                args: [operationId],
                env: [
                  {
                    name: 'RABBITMQ_URL',
                    value:
                      'amqp://guest:guest@rabbitmq.default.svc.cluster.local:5672/',
                  },
                ],
              },
            ],
          },
        },
      },
    };

    try {
      const result = await this.batchV1Api.createNamespacedJob({
        namespace: 'default',
        body: job,
      });
      const jobName = result.metadata?.name;
      if (!jobName) {
        throw new Error('Job name is undefined');
      }

      const timeoutMs = 30000;
      const intervalMs = 2000;
      let elapsed = 0;
      while (elapsed < timeoutMs) {
        try {
          const jobStatus = await this.batchV1Api.readNamespacedJob({
            name: jobName,
            namespace: 'default',
          });
          if (
            jobStatus.status?.active ||
            jobStatus.status?.succeeded ||
            jobStatus.status?.failed
          ) {
            this.logger.log(`Status do job "${jobName}" confirmado.`);
            return;
          }
        } catch (err) {
          this.logger.log(err);
        }
        await new Promise((resolve) => setTimeout(resolve, intervalMs));
        elapsed += intervalMs;
      }
      this.logger.error(
        'Timeout: Job não foi confirmado dentro do tempo esperado.',
      );
      throw new Error('Timeout esperando confirmação da criação do job.');
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : String(error);
      this.logger.error(`Erro ao criar job: ${errorMessage}`);
      throw error;
    }
  }
}
