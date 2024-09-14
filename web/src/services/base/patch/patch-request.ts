import { sendRequest } from "../base";
import { IErrorResponse, IResponse } from "../request-interface";
import { IPatchRequestOption } from "./patch-request-interface";

export default async function patchRequest<T, D>(
	options: IPatchRequestOption<D>,
): Promise<IResponse<T> | IErrorResponse> {
	return sendRequest<T, D>({ method: "PATCH", ...options });
}
