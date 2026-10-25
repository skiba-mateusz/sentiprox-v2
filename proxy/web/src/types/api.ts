export type ApiSuccess<T> = { data: T };
export type ApiFailure = { error: string };

export type PaginatedResult<T> = { items: T; total: number };

export interface LogEntry {
  id: number;
  timestamp: Date;
  ip: string;
  method: string;
  path: string;
  query: string;
  userAgent: string;
  contentType: string;
  route: string;
  source: string;
  category: string;
  score: number;
  action: string;
  statusCode: number;
  latencyMs: number;
}

export interface Summary {
  totalCount: number;
  attacksCount: number;
  topCategory: string;
  topCategoryCount: number;
  totalBlocks: number;
}

export interface Timeserie {
  time: string;
  total: number;
  attacks: number;
}

export interface Category {
  name: string;
  count: number;
}

export interface RequestsParams {
  page: number;
  pageSize: number;
  searchQuery: string;
  method: string;
  since: string;
  to: string;
}
