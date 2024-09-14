import { sendRequest } from "../base";
import { IErrorResponse, IResponse } from "../request-interface";
import { IPutRequestOption } from "./put-request-interface";

export default async function putRequest<T, D>(options: IPutRequestOption<D>): Promise<IResponse<T> | IErrorResponse> {
	return sendRequest<T, D>({ method: "PUT", ...options });
}
