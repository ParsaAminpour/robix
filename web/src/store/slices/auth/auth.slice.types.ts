export interface IUser {
	public_key: string;
	jwt_token: string;
	signed: boolean;
}
export interface AuthStateType {
	user: IUser | null;
}
