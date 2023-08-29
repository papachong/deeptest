import {Action, Mutation} from 'vuex';
import {StoreModuleType} from "@/utils/store";

import {
    clearShareVar,
    quickCreateExtractor,
    getCheckpoint,
    getExtractor,
    getInvocationAsInterface,
    getLastInvocationResp,

    removeCheckpoint,
    removeExtractor,
    removeInvocation,
    removeShareVar,
    saveCheckpoint,
    saveExtractor,
    loadData,
    call,
    save,
    listInvocation,
    listShareVar,
    getSnippet,
    listPostConditions,
    getScript,
    saveScript,
    removeScript,
    createPostConditions,
    removePostConditions,
    movePostConditions,
    disablePostConditions,
    saveAsCase, generateCases,
    getInvocationResult,
    getInvocationLog,
    getPreConditionScript,
    saveResponseDefine, removeCookie, quickCreateCookie, saveCookie, getCookie,
} from './service';
import { serverList, changeServe } from '@/views/project-settings/service';
import {Checkpoint, Cookie, DebugInfo, Extractor, Interface, Response, Script} from "./data";
import {ConditionCategory, ConditionType, UsedBy} from "@/utils/enum";
import {ResponseData} from "@/utils/request";
import {listEnvVarByServer} from "@/services/environment";
import {getResponseKey} from "@/utils/comm";

export interface StateType {
    debugInfo: DebugInfo
    debugData: any;
    invokedMap: any;

    requestData: any;
    responseData: Response;
    consoleData: any[];
    resultData: any[];

    invocationsData: any[];

    preConditions: any[];
    postConditions: any[];
    assertionConditions: any[];
    activeAssertion: any;
    activePostCondition: any;

    extractorData: any;
    checkpointData: any;
    scriptData: any;
    cookieData: any;

    serves: any[];
    currServe: any;
}
const initState: StateType = {
    debugInfo: {} as DebugInfo,
    debugData: {},
    invokedMap: {},

    requestData: {},
    responseData: {} as Response,
    consoleData: [],
    resultData: [],

    invocationsData: [],

    preConditions: [],
    postConditions: [],
    assertionConditions: [],
    activeAssertion: [],
    activePostCondition: [],

    extractorData: {} as Extractor,
    checkpointData: {} as Checkpoint,
    scriptData: {} as Script,
    cookieData: {} as Cookie,

    serves: [],
    currServe: [],
};

export interface ModuleType extends StoreModuleType<StateType> {
    state: StateType;
    mutations: {
        setDebugInfo: Mutation<StateType>;
        setDebugData: Mutation<StateType>;

        clearInvokedMap: Mutation<StateType>;
        putInvokedMap: Mutation<StateType>;

        setRequest: Mutation<StateType>;
        setResponse: Mutation<StateType>;
        setResult: Mutation<StateType>;
        setLog: Mutation<StateType>;

        setInvocations: Mutation<StateType>;
        setServerId: Mutation<StateType>;

        setPostConditions: Mutation<StateType>;
        setAssertionConditions: Mutation<StateType>;
        setActiveAssertion: Mutation<StateType>;
        setActivePostCondition: Mutation<StateType>;

        setExtractor: Mutation<StateType>;
        setCheckpoint: Mutation<StateType>;
        setScript: Mutation<StateType>;
        setScriptContent: Mutation<StateType>;
        setCookie: Mutation<StateType>;

        setPathParams: Mutation<StateType>;
        setShareVars: Mutation<StateType>;
        setEnvVars: Mutation<StateType>;
        setGlobalVars: Mutation<StateType>;

        setUrl: Mutation<StateType>;
        setBaseUrl: Mutation<StateType>;
        setBody: Mutation<StateType>;

        setServes: Mutation<StateType>; // 获取环境列表
        setCurrServe: Mutation<StateType>; // 设置当前所选环境
    };
    actions: {
        loadDataAndInvocations: Action<StateType, StateType>;
        resetDataAndInvocations: Action<StateType, StateType>;
        loadData: Action<StateType, StateType>;
        call: Action<StateType, StateType>;
        refreshInterfaceResultFromScenarioExec: Action<StateType, StateType>;
        save: Action<StateType, StateType>;
        saveAsCase: Action<StateType, StateType>;
        generateCases: Action<StateType, StateType>;

        listInvocation: Action<StateType, StateType>;
        getLastInvocationResp: Action<StateType, StateType>;
        getInvocationResult: Action<StateType, StateType>;
        getInvocationLog: Action<StateType, StateType>;
        getInvocationAsInterface: Action<StateType, StateType>;
        removeInvocation: Action<StateType, StateType>;

        getPreConditionScript: Action<StateType, StateType>;

        listPostCondition: Action<StateType, StateType>;
        listAssertionCondition: Action<StateType, StateType>;
        createPostCondition: Action<StateType, StateType>;
        removePostCondition: Action<StateType, StateType>;
        disablePostCondition: Action<StateType, StateType>;
        movePostCondition: Action<StateType, StateType>;

        getExtractor: Action<StateType, StateType>;
        saveExtractor: Action<StateType, StateType>;
        quickCreateExtractor: Action<StateType, StateType>;
        removeExtractor: Action<StateType, StateType>;

        getCookie: Action<StateType, StateType>;
        saveCookie: Action<StateType, StateType>;
        quickCreateCookie: Action<StateType, StateType>;
        removeCookie: Action<StateType, StateType>;

        getCheckpoint: Action<StateType, StateType>;
        saveCheckpoint: Action<StateType, StateType>;
        removeCheckpoint: Action<StateType, StateType>;

        getScript: Action<StateType, StateType>;
        saveScript: Action<StateType, StateType>;
        removeScript: Action<StateType, StateType>;
        addSnippet: Action<StateType, StateType>;

        listShareVar: Action<StateType, StateType>;
        removeShareVar: Action<StateType, StateType>;
        clearShareVar: Action<StateType, StateType>;

        updateUrl: Action<StateType, StateType>;
        updateBaseUrl:Action<StateType, StateType>;
        updateBody: Action<StateType, StateType>;

        changeServer: Action<StateType, StateType>;

        saveResponseDefine:Action<StateType, StateType>

        listServes: Action<StateType, StateType>;
    };
}

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'Debug',
    state: {
        ...initState
    },
    mutations: {
        setDebugInfo(state, payload) {
            state.debugInfo = payload;
        },
        setServerId(state, id) {
            state.debugData.serverId = id;
        },

        setDebugData(state, payload) {
            state.debugData = payload;
        },
        clearInvokedMap(state) {
            state.invokedMap = {}
        },
        putInvokedMap(state) {
            const key = getResponseKey(state.debugInfo)
            state.invokedMap[key] = true;
        },
        setRequest(state, payload) {
            state.requestData = payload;
        },
        setResponse(state, payload) {
            state.responseData = payload;
        },
        setResult(state, payload) {
            state.resultData = payload;
        },
        setLog(state, payload) {
            state.consoleData = payload;
        },

        setInvocations(state, payload) {
            state.invocationsData = payload;
        },

        setPostConditions(state, payload) {
            state.postConditions = payload;
        },
        setAssertionConditions(state, payload) {
            state.assertionConditions = payload;
        },

        setActiveAssertion(state, payload) {
            if (state.activeAssertion.id === payload.id) {
                state.activeAssertion = {}
            } else {
                state.activeAssertion = payload;
            }
        },
        setActivePostCondition(state, payload) {
            if (state.activePostCondition.id === payload.id) {
                state.activePostCondition = {}
            } else {
                state.activePostCondition = payload;
            }
        },

        setExtractor(state, payload) {
            state.extractorData = payload;
        },
        setCheckpoint(state, payload) {
            state.checkpointData = payload;
        },
        setScript(state, payload) {
            state.scriptData = payload;
        },
        setCookie(state, payload) {
            state.cookieData = payload;
        },
        setScriptContent(state, content) {
            console.log('setScriptContent', content)
            state.scriptData.content = content;
        },

        setPathParams(state, payload) {
            console.log('set debugData pathParams')
            state.debugData.pathParams = payload;
        },
        setShareVars(state, payload) {
            console.log('set debugData shareVars')
            state.debugData.shareVars = payload;
        },
        setEnvVars(state, payload) {
            console.log('set debugData envVars')
            state.debugData.envVars = payload;
        },
        setGlobalVars(state, payload) {
            console.log('set debugData globalVars')
            state.debugData.globalVars = payload;
        },

        setBaseUrl(state, payload) {
            console.log('set debugData baseUrl')
            state.debugData.baseUrl = payload;
        },
        setUrl(state, payload) {
            console.log('set debugData url')
            state.debugData.url = payload;
        },
        setBody(state, payload) {
            console.log('set debugData body')
            state.debugData.body = payload;
        },
        setServes(state, payload) {
            console.log('292---- ', payload);
            state.serves = payload || [];
        },
        setCurrServe(state, payload) {
            state.currServe = payload;
        }
    },
    actions: {
        // debug
        async loadDataAndInvocations({commit, dispatch, state}, data) {
            try {
                await dispatch('loadData', data)

                dispatch('listInvocation', {
                    debugInterfaceId: state.debugInfo.debugInterfaceId,
                    endpointInterfaceId: state.debugInfo.endpointInterfaceId,
                    diagnoseInterfaceId: state.debugInfo.diagnoseInterfaceId,
                    caseInterfaceId: state.debugInfo.caseInterfaceId,
                })

                return true;
            } catch (error) {
                return false;
            }
        },
        async loadData({commit, state, dispatch}, data) {
            try {
                const resp: ResponseData = await loadData(data);
                // console.log('666666', resp.data.url)
                if (resp.code != 0) return false;

                await commit('setDebugInfo', {
                    debugInterfaceId: resp.data.debugInterfaceId,
                    endpointInterfaceId: data.endpointInterfaceId,
                    scenarioProcessorId  : data.scenarioProcessorId,
                    diagnoseInterfaceId  : data.diagnoseInterfaceId,
                    caseInterfaceId: data.caseInterfaceId,
                    usedBy:          data.usedBy,
                } as DebugInfo);

                await commit('setDebugData', resp.data);

                const key = getResponseKey(state.debugInfo)
                if (state.invokedMap[key]) {
                    await dispatch('getLastInvocationResp')
                }

                return true;
            } catch (error) {
                return false;
            }
        },
        async save({commit}, payload: any) {
            const resp = await  save(payload)
            if (resp.code === 0) {
                commit('setDebugData', { ...resp.data, action: 'save' });

                return true;
            } else {
                return false
            }
        },
        async saveAsCase({commit}, payload: any) {
            const resp = await  saveAsCase(payload)
            if (resp.code === 0) {
                // commit('setDebugData', resp.data);
                return true;
            } else {
                return false
            }
        },
        async generateCases({commit}, payload: any) {
            const resp = await  generateCases(payload)
            if (resp.code === 0) {
                return true;
            } else {
                return false
            }
        },
        async resetDataAndInvocations({commit, dispatch, state}) {
            commit('setDebugInfo', {});
            commit('setDebugData', {});
            commit('setRequest', {});
            commit('setResponse', {});
            commit('setResult', []);
            commit('setLog', []);
            commit('setInvocations', []);
        },

        async call({commit, dispatch, state}, data: any) {
            console.log('call')

            commit('setRequest', {});
            commit('setResponse', {});
            const response = await call(data)

            if (response.code === 0) {
                commit('setRequest', response.data.req);
                commit('setResponse', response.data.resp);

                await dispatch('listInvocation')
                await dispatch('getLastInvocationResp')

                commit('putInvokedMap')

                await dispatch('listShareVar');

                await dispatch('listPostCondition');
                await dispatch('listAssertionCondition');

                return true;
            } else {
                return false
            }
        },

        async refreshInterfaceResultFromScenarioExec({commit, dispatch, state}) {
            await dispatch('getLastInvocationResp')
            await commit('putInvokedMap')
        },

        // invocation
        async listInvocation({commit, state}) {
            try {
                const resp = await listInvocation({
                    debugInterfaceId: state.debugInfo.debugInterfaceId,
                    endpointInterfaceId: state.debugInfo.endpointInterfaceId,
                    diagnoseInterfaceId: state.debugInfo.diagnoseInterfaceId,
                    caseInterfaceId: state.debugInfo.caseInterfaceId,
                } as DebugInfo);
                const {data} = resp;
                commit('setInvocations', data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async getLastInvocationResp({commit, dispatch, state}) {
            const response = await getLastInvocationResp( {
                debugInterfaceId: state.debugInfo.debugInterfaceId,
                endpointInterfaceId: state.debugInfo.endpointInterfaceId,
                diagnoseInterfaceId: state.debugInfo.diagnoseInterfaceId,
                caseInterfaceId: state.debugInfo.caseInterfaceId,
            } as DebugInfo);

            commit('setRequest', response.data.req);
            commit('setResponse', response.data.resp);
            return true;
        },
        async getInvocationResult({commit, dispatch, state}, invokeId: number) {
            const response = await getInvocationResult(invokeId);
            commit('setResult', response.data);
            return true;
        },
        async getInvocationLog({commit, dispatch, state}, invokeId: number) {
            const response = await getInvocationLog(invokeId);
            commit('setLog', response.data);
            return true;
        },

        async getInvocationAsInterface({commit}, id: number) {
            try {
                const resp = await getInvocationAsInterface(id);
                const {data} = resp;

                commit('setDebugData', data.debugData);
                commit('setRequest', data.req);
                commit('setResponse', data.resp);
                return true;
            } catch (error) {
                return false;
            }
        },
        async removeInvocation({commit, dispatch, state}, id: number) {
            try {
                await removeInvocation(id);
                dispatch('listInvocation', {
                    endpointInterfaceId: state.debugInfo.endpointInterfaceId,
                });
                return true;
            } catch (error) {
                return false;
            }
        },

        // conditions
        async getPreConditionScript({commit, state}) {
            try {
                const resp = await getPreConditionScript(state.debugInfo.debugInterfaceId, state.debugData.endpointInterfaceId);
                const {data} = resp;
                commit('setScript', data);
                return true;
            } catch (error) {
                return false;
            }
        },

        async listPostCondition({commit, state}) {
            try {
                const resp = await listPostConditions(state.debugInfo.debugInterfaceId, state.debugData.endpointInterfaceId,
                    ConditionCategory.console);
                const {data} = resp;
                commit('setPostConditions', data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async listAssertionCondition({commit, state}) {
            try {
                const resp = await listPostConditions(state.debugInfo.debugInterfaceId, state.debugData.endpointInterfaceId,
                    ConditionCategory.assert);

                const {data} = resp;
                commit('setAssertionConditions', data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async createPostCondition({commit, dispatch, state}, payload: any) {
            try {
                await createPostConditions(payload);

                if (payload.entityType === ConditionType.checkpoint) {
                    await dispatch('listAssertionCondition');

                    const len = state.assertionConditions.length
                    if (len > 0) {
                        commit('setActiveAssertion', state.assertionConditions[len-1]);
                    }

                } else {
                    await dispatch('listPostCondition');

                    const len = state.postConditions.length
                    if (len > 0) {
                        commit('setActivePostCondition', state.postConditions[len-1]);
                    }
                }
                return true;
            } catch (error) {
                return false;
            }
        },
        async disablePostCondition({commit, dispatch, state}, payload: any) {
            try {
                await disablePostConditions(payload.id);
                if (payload.entityType === ConditionType.checkpoint) {
                    dispatch('listAssertionCondition');
                } else {
                    dispatch('listPostCondition');
                }
                return true;
            } catch (error) {
                return false;
            }
        },
        async removePostCondition({commit, dispatch, state}, payload: any) {
            try {
                await removePostConditions(payload.id);
                if (payload.entityType === ConditionType.checkpoint) {
                    dispatch('listAssertionCondition');
                } else {
                    dispatch('listPostCondition');
                }
                return true;
            } catch (error) {
                return false;
            }
        },
        async movePostCondition({commit, dispatch, state}, payload: any) {
            try {
                await movePostConditions(payload);
                if (payload.entityType === ConditionType.checkpoint) {
                    dispatch('listAssertionCondition');
                } else {
                    dispatch('listPostCondition');
                }
                return true;
            } catch (error) {
                return false;
            }
        },

        // extractor
        async getExtractor({commit}, id: number) {
            try {
                const response = await getExtractor(id);
                const {data} = response;

                commit('setExtractor', data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async saveExtractor({commit, dispatch, state}, payload: any) {
            try {
                await saveExtractor(payload);
                dispatch('listPostCondition');
                return true;
            } catch (error) {
                return false;
            }
        },
        async quickCreateExtractor({commit, dispatch, state}, payload: any) {
            try {
                await quickCreateExtractor(payload);
                dispatch('listPostCondition');
                dispatch('listShareVar');
                return true;
            } catch (error) {
                return false;
            }
        },
        async removeExtractor({commit, dispatch, state}, payload) {
            try {
                await removeExtractor(payload.id);

                dispatch('listPostCondition');
                return true;
            } catch (error) {
                return false;
            }
        },

        // cookie
        async getCookie({commit}, id: number) {
            try {
                const response = await getCookie(id);
                const {data} = response;

                commit('setCookie', data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async saveCookie({commit, dispatch, state}, payload: any) {
            try {
                await saveCookie(payload);
                dispatch('listPostCondition');
                return true;
            } catch (error) {
                return false;
            }
        },
        async quickCreateCookie({commit, dispatch, state}, payload: any) {
            try {
                await quickCreateCookie(payload);
                dispatch('listPostCondition');
                return true;
            } catch (error) {
                return false;
            }
        },
        async removeCookie({commit, dispatch, state}, payload) {
            try {
                await removeCookie(payload.id);

                dispatch('listPostCondition');
                return true;
            } catch (error) {
                return false;
            }
        },

        // checkpoint
        async getCheckpoint({commit}, id: number) {
            try {
                const response = await getCheckpoint(id);
                const {data} = response;

                commit('setCheckpoint', data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async saveCheckpoint({commit, dispatch, state}, payload: any) {
            try {
                await saveCheckpoint(payload);
                dispatch('listPostCondition', UsedBy.InterfaceDebug);
                return true
            } catch (error) {
                return false;
            }
        },
        async removeCheckpoint({commit, dispatch, state}, id: number) {
            try {
                await removeCheckpoint(id);

                dispatch('listPostCondition', UsedBy.InterfaceDebug);
                return true;
            } catch (error) {
                return false;
            }
        },

        // script
        async getScript({commit}, id: number) {
            try {
                const response = await getScript(id);
                const {data} = response;

                commit('setScript', data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async saveScript({commit, dispatch, state}, payload: any) {
            try {
                await saveScript(payload);
                dispatch('listPostCondition', UsedBy.InterfaceDebug);
                return true
            } catch (error) {
                return false;
            }
        },
        async removeScript({commit, dispatch, state}, id: number) {
            try {
                await removeScript(id);

                dispatch('listPostCondition', UsedBy.InterfaceDebug);
                return true;
            } catch (error) {
                return false;
            }
        },
        async addSnippet({commit, dispatch, state}, name: string) {
            const json = await getSnippet(name)
            if (json.code === 0) {
                let script = (state.scriptData.content ? state.scriptData.content: '') + '\n' +  json.data.script
                script = script.trim()
                commit('setScriptContent', script);
            }

            return true;
        },

        // shared variable
        async listShareVar({commit, dispatch, state}) {
            try {
                const resp = await listShareVar(state.debugInfo, UsedBy.InterfaceDebug);
                const {data} = resp;
                commit('setShareVars', data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async clearShareVar({commit, dispatch, state}, payload: any) {
            try {
                console.log('debugInfo', state.debugInfo)

                await clearShareVar(state.debugInfo);
                dispatch('listShareVar');

                return true;
            } catch (error) {
                return false;
            }
        },
        async removeShareVar({commit, dispatch, state}, payload: any) {
            try {
                const resp = await removeShareVar(payload.id);
                const {data} = resp;
                dispatch('listShareVar');

                return true;
            } catch (error) {
                return false;
            }
        },
        async updateUrl({commit, dispatch, state}, body: string) {
            commit('setUrl', body);
            return true;
        },
        async updateBaseUrl({commit, dispatch, state}, body: string) {
            commit('setBaseUrl', body);
            return true;
        },
        async updateBody({commit, dispatch, state}, body: string) {
            commit('setBody', body);
            return true;
        },

        async changeServer({commit, dispatch, state}, payload: {serverId: number, requestEnvVars: boolean}) {
            const { serverId, requestEnvVars = true } = payload;
            const res = await changeServe({ serverId });
            if (res.code === 0) {
                commit('setCurrServe', res.data);
            }
            if (requestEnvVars) {
                const json = await listEnvVarByServer(serverId)
                if (json.code === 0) {
                    commit('setServerId', serverId);
                    commit('setEnvVars', json.data);
                }
            }
            return true;
        },

        async saveResponseDefine({commit, dispatch, state}, payload: any) {
            try {
                await saveResponseDefine(payload);
                state.debugData.responseDefine.entityData.disabled = payload.disabled
                state.debugData.responseDefine.entityData.code = payload.code
                return true
            } catch (error) {
                return false;
            }
        },
        async listServes({ commit }, payload: { serveId: number }) {
            const res = await serverList(payload);
            if (res.code === 0) {
                const servers = (res.data.servers || []).map((item: any) => {
                    item.label = item.environmentName;
                    item.value = item.environmentId;
                    return item;
                })
                commit('setServes', servers);
                commit('setCurrServe', res.data.currServer);
            }
        }
    }
};

export default StoreModel;
