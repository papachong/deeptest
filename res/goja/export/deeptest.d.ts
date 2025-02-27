declare type Param = {
    name: string;
    value: string;
    paramIn?: string;
    disabled?: boolean;
    type?: string;
}
declare type Header = {
    name: string;
    value: string;
    disabled?:    boolean;
    type?: string;
}
declare type ExecCookie = {
    name: string;
    value: any;
    path?: string;

    domain?: string;
    expireTime?: Date;
}
declare type BodyFormDataItem = {
    name: string;
    value: string;
    type: string;
}

declare type BodyFormUrlEncodedItem = {
    name: string;
    value: string;
}

declare type Request = {
    method: string;
    url: string;
    queryParams: Param[];
    pathParams: Param[];
    headers: Header[];
    cookies: ExecCookie[];

    body: string;
    bodyFormData:       BodyFormDataItem[];
    bodyFormUrlencoded: BodyFormUrlEncodedItem[];
    bodyType: string;
};
declare type Response = {
    statusCode: number;
    statusContent: string;

    headers: Header[];
    cookies: ExecCookie[];

    data: any;
    contentType: string;

    contentCharset: string;
    contentLength: number;
}
declare type ResponseWithContent = {
    statusCode: number;
    statusContent: string;

    headers: Header[];
    cookies: ExecCookie[];

    content: string;
    contentType: string;

    contentCharset: string;
    contentLength: number;
}

declare global {
    const dt: {
        datapool: {
            get: (datapool_name: string, variable_name: string, seq: string) => any,
        },
        variables: {
            get: (variable_name: string) => any,
            set: (variable_name: string, variable_value: any) => {},
            clear: (variable_name: string) => {},
        },

        request: Request,
        response: Response,

        sendRequest: (urlOrConfig: string | object, callback: (error, response: ResponseWithContent) => void) => void,
    }

    const log : (obj: any) => {}
}

export default {};
