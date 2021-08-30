export interface alert {
  id: string;
  title: string;
  description: string;
  detector: string;
  status: string;
  severity: string;
  created_at: number;
  detected_at: number;
  closed_at: number;

  attributes: attribute[];
  task_logs: taskLog[];
  references: reference[];
}

export interface attribute {
  id: number;
  key: string;
  value: string;
  type: string;
  context: string[];

  annotations: annotation[];
  actions: action[];
}

export interface action {
  id: string;
  name: string;
}

export interface annotation {
  id: number;
  timestamp: number;
  source: string;
  name: string;
  value: string;
}

export interface taskLog {
  id: number;
  task_name: string;
  optional: boolean;
  stage: number;
  started_at: number;
  exited_at: number;
  log: string;
  errmsg: string;
  err_values: string[];
  stack_trace: string[];
  status: string;
}

export interface reference {
  id: string;
  source: string;
  title: string;
  url: string;
  comment: string;
}
