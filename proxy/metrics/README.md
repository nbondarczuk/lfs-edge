# Server Prometheus metrics

## General metrics

- ReportLatencyMetric:

  It reports the latency of the specified operation to the
  specified summary vector metric. The label is used to partition the resulting
  histogram.

## Rest metrics

- MetricRestLatency:

  REST request processing latency is partitioned by the REST method. It uses
  custom buckets based on the expected request duration.

- MetricRequestCount:

  Number of REST requests received.

- MetricGetFileInternalErrors:

  Number of internal errors encountered when processing get file requests.

- MetricGetFileBadRequests:

  Number of bad get file requests encountered.

- MetricGetFileNotFoundErrors:

  Number of get file requests where the requested file was not found.

- MetricGetFileResponses:

  Number of successful get file requests served.
