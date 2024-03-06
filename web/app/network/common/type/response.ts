export interface ResponseError {
  status: number;
  error: string;
}

export interface ResponseSuccess<T> {
  data: T;
}
