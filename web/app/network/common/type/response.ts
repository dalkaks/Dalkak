export interface ResponseError {
  status: number;
  error: {
    message: string;
  };
}

export interface ResponseSuccess<T> {
  data: T;
}
