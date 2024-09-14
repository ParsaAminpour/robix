import { sendRequest } from "../base";
import { IErrorResponse, IResponse } from "../request-interface";
import { IDeleteRequestOption } from "./delete-request-interface";

export default async function deleteRequest<T, D>(
	options: IDeleteRequestOption<D>,
): Promise<IResponse<T> | IErrorResponse> {
	return sendRequest<T, D>({ method: "DELETE", ...options });
}
