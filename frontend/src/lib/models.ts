export class Agent {
    id: number
    username: string
    email: string
    role: number
    mfa: boolean

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source)
        this.id = source['id']
        this.username = source['username']
        this.email = source['email']
        this.role = source['role']
        this.mfa = source['two_factor_auth']['verified']
    }

    static createFrom(source: any = {}) {
        return new Agent(source)
    }
}

export interface Result<T> {
    totalRecords: number,
    totalPage: number,
    offset: number,
    page: number,
    prevPage: number,
    nextPage: number,
    results: Array<T>
}

export interface TokenInfo {
    token: string,
    expire: string,
    orig_iat: string
}