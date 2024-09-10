import { PropsWithChildren } from "react";

export interface BaseProps<T = unknown> extends React.FC<PropsWithChildren<T>> {}
