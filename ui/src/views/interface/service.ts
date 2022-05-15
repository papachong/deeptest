import request from '@/utils/request';
import {Interface, OAuth20} from "@/views/interface/data";
import {isInArray} from "@/utils/array";
import {CheckpointOperator} from "@/views/interface/consts";

const apiPath = 'interfaces';
const apiVocation = 'invocations';
const apiAuth = 'auth';
const apiEnvironment = 'environments'
const apiEnvironmentVar = `${apiEnvironment}/vars`

const apiExtractor = 'extractors'
const apiCheckpoint = 'checkpoints'

// interface
export async function saveInterface(interf: Interface): Promise<any> {
    return request({
        url: `/${apiPath}/saveInterface`,
        method: 'post',
        data: interf,
    });
}

export async function load(): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
    });
}

export async function get(id: number): Promise<any> {
    return request({url: `/${apiPath}/${id}`});
}

export async function create(data): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'POST',
        data: data,
    });
}

export async function update(id: number, params: any): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'PUT',
        data: params,
    });
}

export async function updateNameReq(id: number, name: string): Promise<any> {
    const data = {id: id, name: name}
    return request({
        url: `/${apiPath}/updateName`,
        method: 'PUT',
        data: data,
    });
}

export async function remove(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
    });
}

export async function move(data: any): Promise<any> {
    return request({
        url: `/${apiPath}/move`,
        method: 'post',
        data: data,
    });
}

// invocation
export async function invoke(interf: Interface): Promise<any> {
    return request({
        url: `/${apiVocation}/invoke`,
        method: 'post',
        data: interf,
    });
}
export async function listInvocation(interfaceId: number): Promise<any> {
    const params = {interfaceId: interfaceId}

    return request({
        url: `/${apiVocation}`,
        method: 'GET',
        params,
    });
}
export async function getInvocationAsInterface(id: number): Promise<any> {
    return request({url: `/${apiVocation}/${id}`});
}
export async function removeInvocation(id: number): Promise<any> {
    return request({
        url: `/${apiVocation}/${id}`,
        method: 'DELETE',
    });
}

// auth
export async function genOAuth2AccessToken(oauth: OAuth20): Promise<any> {
    return request({
        url: `/${apiAuth}/oauth2Authorization`,
        method: 'post',
        data: oauth,
    });
}

// environment
export async function listEnvironment(): Promise<any> {
    const params = {}

    return request({
        url: `/${apiEnvironment}`,
        method: 'GET',
        params,
    });
}
export async function getEnvironment(id: number, interfaceId: number): Promise<any> {
    const params = {interfaceId: interfaceId}
    return request({
        url: `/${apiEnvironment}/${id}`,
        method: 'GET',
        params
    });
}

export async function changeEnvironment(id, interfaceId): Promise<any> {
    const params = {id, interfaceId}

    return request({
        url: `/${apiEnvironment}/changeEnvironment`,
        method: 'POST',
        params,
    });
}
export async function saveEnvironment(data): Promise<any> {
    return request({
        url: `/${apiEnvironment}`,
        method: data.id ? 'PUT' : 'POST',
        data: data,
    });
}
export async function copyEnvironment(id): Promise<any> {
    const params = {id: id}
    return request({
        url: `/${apiEnvironment}/copyEnvironment`,
        method: 'POST',
        params,
    });
}
export async function removeEnvironment(id: number): Promise<any> {
    return request({
        url: `/${apiEnvironment}/${id}`,
        method: 'DELETE',
    });
}

// environment var
export async function saveEnvironmentVar(data): Promise<any> {
    return request({
        url: `/${apiEnvironmentVar}`,
        method: data.id ? 'PUT' : 'POST',
        data: data,
    });
}
export async function removeEnvironmentVar(id: number): Promise<any> {
    return request({
        url: `/${apiEnvironmentVar}/${id}`,
        method: 'DELETE',
    });
}
export async function clearEnvironmentVar(environmentId: number): Promise<any> {
    const params = {environmentId: environmentId}
    return request({
        url: `/${apiEnvironmentVar}/clear`,
        method: 'POST',
        params,
    });
}

// helper
export function prepareDataForRequest(data: any) {
    data.params = data.params.filter((item) => {
        return !!item.name
    })
    data.headers = data.headers.filter((item) => {
        return !!item.name
    })

    return data
}

export function getNodeMap(treeNode: any, mp: any): void {
    if (!treeNode) return

    mp[treeNode.id] = treeNode

    if (treeNode.children) {
        treeNode.children.forEach((item, index) => {
            getNodeMap(item, mp)
        })
    }

    return
}

export function expandAllKeys(treeMap: any, isExpand: boolean): number[] {
    const keys = new Array<number>()
    if (!isExpand) return keys

    Object.keys(treeMap).forEach(key => {
        if (!keys.includes(+key)) keys.push(+key)
    })

    return keys
}

export function expandOneKey(treeMap: any, key: number, expandedKeys: number[]) {
    if (!expandedKeys.includes(key)) expandedKeys.push(key)

    const parentId = treeMap[key].parentId
    if (parentId) {
        expandOneKey(treeMap, parentId, expandedKeys)
    }
}

export function getCodeLangByType(type) {
    type = type.split('/')[1]

    if (isInArray(type, ['json', 'xml', 'html'])) {
        return type
    } else {
        return 'plaintext'
    }
}

export const getEnumSelectItems = (enumDef) => {
    const arr : any[] = []

    for (const item in enumDef) {
        console.log('======' , item, enumDef[item])
        arr.push({label: enumDef[item], value: item})
    }

    return arr
}
export const getOperatorsForString = () => {
    const arr : any[] = []

    arr.push({label: CheckpointOperator.equal, value: CheckpointOperator[CheckpointOperator.equal]})
    arr.push({label: CheckpointOperator.notEqual, value: CheckpointOperator[CheckpointOperator.notEqual]})
    arr.push({label: CheckpointOperator.contain, value: CheckpointOperator[CheckpointOperator.contain]})

    return arr
}

export const getOperatorsForCode = () => {
    const arr : any[] = []

    arr.push({label: CheckpointOperator.equal, value: CheckpointOperator[CheckpointOperator.equal]})
    arr.push({label: CheckpointOperator.notEqual, value: CheckpointOperator[CheckpointOperator.notEqual]})

    return arr
}

// extractor
export async function listExtractor(interfaceId: number): Promise<any> {
    const params = {interfaceId}

    return request({
        url: `/${apiExtractor}`,
        method: 'GET',
        params,
    });
}
export async function getExtractor(id: number): Promise<any> {
    return request({
        url: `/${apiExtractor}/${id}`,
        method: 'GET',
    });
}
export async function saveExtractor(data): Promise<any> {
    return request({
        url: `/${apiExtractor}`,
        method: data.id ? 'PUT' : 'POST',
        data: data,
    });
}
export async function removeExtractor(id: number): Promise<any> {
    return request({
        url: `/${apiExtractor}/${id}`,
        method: 'DELETE',
    });
}
export async function listExtractorVariable(interfaceId: number): Promise<any> {
    const params = {interfaceId}
    return request({
        url: `/${apiExtractor}/listExtractorVariable`,
        method: 'GET',
        params,
    });
}

// checkpoint
export async function listCheckpoint(interfaceId: number): Promise<any> {
    const params = {interfaceId}

    return request({
        url: `/${apiCheckpoint}`,
        method: 'GET',
        params,
    });
}
export async function getCheckpoint(id: number): Promise<any> {
    return request({
        url: `/${apiCheckpoint}/${id}`,
        method: 'GET',
    });
}
export async function saveCheckpoint(data): Promise<any> {
    return request({
        url: `/${apiCheckpoint}`,
        method: data.id ? 'PUT' : 'POST',
        data: data,
    });
}
export async function removeCheckpoint(id: number): Promise<any> {
    return request({
        url: `/${apiCheckpoint}/${id}`,
        method: 'DELETE',
    });
}