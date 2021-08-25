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
}

export interface attribute {
  id: string;
  key: string;
  value: string;
  type: string;
  context: string[];

  annotations: annotation[];
}

export interface annotation {
  timestamp: number;
  source: string;
  name: string;
  value: string;
}
