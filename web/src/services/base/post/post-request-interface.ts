import { AxiosRequestConfig } from "axios";

export interface IPostRequestOption<D> extends Omit<AxiosRequestConfig<D>, "method"> {}
