export default class Client {
    notes: notes.ServiceClient

    constructor(environment: string = "prod", token?: string) {
        const base = new BaseClient(environment, token)
        this.notes = new notes.ServiceClient(base)
    }
}

export namespace notes {
    export interface ListResponse {
        notes: Note[]
    }

    export interface NewNoteRequest {
        front: string
        back: string
    }

    export interface Note {
        id: string
        front: string
        back: string
        next_review: string
    }

    export interface ReviewNoteRequest {
        answer: string
    }

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        public ListNotes(): Promise<ListResponse> {
            return this.baseClient.do<ListResponse>("GET", `/note`)
        }

        public NewNote(params: NewNoteRequest): Promise<Note> {
            return this.baseClient.do<Note>("POST", `/note`, params)
        }

        public ReviewNote(id: string, params: ReviewNoteRequest): Promise<Note> {
            return this.baseClient.do<Note>("POST", `/note/${id}`, params)
        }
    }
}

class BaseClient {
    baseURL: string
    headers: {[key: string]: string}

    constructor(environment: string, token?: string) {
        this.headers = {"Content-Type": "application/json"}
        if (token !== undefined) {
            this.headers["Authorization"] = "Bearer " + token
        }
        if (environment === "local") {
            this.baseURL = "http://localhost:4000"
        } else {
            this.baseURL = `https://encore-spaced-repetition-im32.encoreapi.com/${environment}`
        }
    }

    public async do<T>(method: string, path: string, req?: any): Promise<T> {
        let response = await fetch(this.baseURL + path, {
            method: method,
            headers: this.headers,
            body: req !== undefined ? JSON.stringify(req) : undefined
        })
        if (!response.ok) {
            let body = await response.text()
            throw new Error("request failed: " + body)
        }
        return <T>(await response.json())
    }

    public async doVoid(method: string, path: string, req?: any): Promise<void> {
        let response = await fetch(this.baseURL + path, {
            method: method,
            headers: this.headers,
            body: req !== undefined ? JSON.stringify(req) : undefined
        })
        if (!response.ok) {
            let body = await response.text()
            throw new Error("request failed: " + body)
        }
        await response.text()
    }
}

function encodeQuery(parts: any[]): string {
    const pairs = []
    for (let i = 0; i < parts.length; i += 2) {
        const key = parts[i]
        let val = parts[i+1]
        if (!Array.isArray(val)) {
            val = [val]
        }
        for (const v of val) {
            pairs.push(`${key}=${encodeURIComponent(v)}`)
        }
    }
    return pairs.join("&")
}
