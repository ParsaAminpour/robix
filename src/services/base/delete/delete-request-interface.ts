import { AxiosRequestConfig } from "axios";

export interface IDeleteRequestOption<D> extends Omit<AxiosRequestConfig<D>, "method"> {}
