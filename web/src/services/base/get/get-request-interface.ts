import { AxiosRequestConfig } from "axios";

export interface IGetRequestOption extends Omit<AxiosRequestConfig, "method" | "data"> {}
