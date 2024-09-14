import { AxiosRequestConfig } from "axios";

export interface IPatchRequestOption<D> extends Omit<AxiosRequestConfig<D>, "method"> {}
