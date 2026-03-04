export const SHARE_CONSTANTS = {
    TOKEN_LENGTH: 32,
    DEFAULT_EXPIRES_HOURS: 0,
    DEFAULT_RULE_URL: 'https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/config/ACL4SSR_Online_Full_NoAuto.ini',
    DEFAULT_RENAME_TEMPLATE: '{{.Country.Emoji}}{{.Country.NameZh}} {{.Delay}} {{.Count}}',
} as const

export const SUBSCRIPTION_TARGETS = [
    { value: 'clash', label: 'Clash/Mihomo' },
    { value: 'surge&ver=4', label: 'Surge4/5' },
    { value: 'singbox', label: 'Sing-Box' },
    { value: 'v2ray', label: 'V2Ray' },
    { value: 'trojan', label: 'Trojan' },
    { value: 'ssr', label: 'ShadowsocksR' },
    { value: 'mixed', label: '混合订阅（mixed）' },
    { value: 'surfboard', label: 'Surfboard' },
    { value: 'quan', label: 'Quantumult' },
    { value: 'quanx', label: 'Quantumult X' },
    { value: 'loon', label: 'Loon' },
    { value: 'mellow', label: 'Mellow' },
    { value: 'surge&ver=3', label: 'Surge3' },
    { value: 'surge&ver=2', label: 'Surge2' },
    { value: 'clashr', label: 'ClashR' },
    { value: 'ss', label: 'Shadowsocks(SIP002)' },
    { value: 'sssub', label: 'Shadowsocks Android(SIP008)' },
    { value: 'ssd', label: 'ShadowsocksD' },
    { value: 'auto', label: '自动判断客户端' },
] as const

export const FORM_VALIDATION = {
    NAME_REQUIRED: '请输入分享名称',
    POSITIVE_NUMBER: '请输入正数',
    VALID_COUNTRY_CODE: '请输入有效的国家代码',
} as const

export const UI_TEXT = {
    CREATE_SHARE: '创建分享',
    EDIT_SHARE: '编辑分享',
    UPDATE: '更新',
    CREATE: '创建',
    CANCEL: '取消',
    DELETE: '删除',
    COPY: '复制',
    LOADING: '加载中...',
    NO_DATA: '暂无数据',
    CONFIRM_DELETE: '确认删除',
    DELETE_CONFIRM_MESSAGE: '您确定要删除分享 "{name}" 吗？此操作无法撤销。',
    COPY_SUCCESS: '复制成功',
    COPY_FAILED: '复制失败',
    CREATE_SUCCESS: '分享创建成功',
    UPDATE_SUCCESS: '分享更新成功',
    DELETE_SUCCESS: '分享删除成功',
    CREATE_FAILED: '创建分享失败',
    UPDATE_FAILED: '更新分享失败',
    DELETE_FAILED: '删除分享失败',
} as const