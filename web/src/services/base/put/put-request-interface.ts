import { AxiosRequestConfig } from "axios";

export interface IPutRequestOption<D> extends Omit<AxiosRequestConfig<D>, "method"> {}
