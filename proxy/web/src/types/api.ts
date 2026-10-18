export type ApiSuccess<T> = { data: T };
export type ApiFailure = { error: string };

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
