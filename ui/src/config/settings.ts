import { RoutesDataItem } from "@/utils/routes";
import {getShowRightBar} from "@/utils/cache";

/**
 * 站点配置
 * @author LiQingSong
 */
export interface SettingsType {
    /**
     * 站点名称
     */
    siteTitle: string;

    /**
     * 顶部菜单开启
     */
    topNavEnable: boolean;

    /**
     * 头部固定开启
     */
    headFixed: boolean;

    /**
     * tab菜单开启
     */
     tabNavEnable: boolean;

     /**
     * 站点首页路由
     */
    homeRouteItem: RoutesDataItem;

    /**
     * 站点本地存储Token 的 Key值
     */
    siteTokenKey: string;

    /**
     * 站点本地存储当前项目 的 Key值
     */
    currProjectId: string;

    settings: string;
    showRightBar: string;
    expandedKeys: string;

    eventNotify,
    eventWebSocketConnStatus: string,
    eventWebSocketMsg: string,
    eventEditorContainerHeightChanged: string,
    webSocketRoom: string,
    electronMsg: string,
    electronMsgReplay: string,

    /**
     * Ajax请求头发送Token 的 Key值
     */
    ajaxHeadersTokenKey: string;

    /**
     * Ajax返回值不参加统一验证的api地址
     */
    ajaxResponseNoVerifyUrl: string[];

    /**
     * iconfont.cn 项目在线生成的 js 地址
     */
    iconfontUrl: string[];
}

const settings: SettingsType = {
    siteTitle: 'ADMIN-ANTD-VUE',
    topNavEnable: true,
    headFixed: true,
    tabNavEnable: false,
    homeRouteItem: {
        icon: 'project',
        title: 'index-layout.menu.project.list',
        path: '/project/index',
        component: ()=> import('@/views/home/index.vue')
    },
    siteTokenKey: 'admin_antd_vue_token',
    currProjectId: 'curr_project_id',

    settings: 'settings',
    showRightBar: 'showRightBar',
    expandedKeys: 'deeptest-expandedKeys',

    eventNotify: 'eventNotify',
    eventWebSocketConnStatus: 'eventWebSocketStatus',
    eventWebSocketMsg: 'eventWebSocketMsg',
    eventEditorContainerHeightChanged: 'eventWebSocketMsg',
    webSocketRoom: 'webSocketRoom',
    electronMsg: 'electronMsg',
    electronMsgReplay: 'electronMsgReplay',

    // ajaxHeadersTokenKey: 'x-token',
    ajaxHeadersTokenKey: 'Authorization',
    ajaxResponseNoVerifyUrl: [
        '/user/login', // 用户登录
    ],
    iconfontUrl: [],
};

export default settings;
