export interface GenConfig {
    filter: ShareFilter
    rename: string
    proxy: boolean
    sub_converter: SubConverterConfig
}

export interface SubConverterConfig {
    target: string
    config: string
}
export interface ShareRequest {
    enable: boolean
    name: string
    token: string
    gen: GenConfig
    max_access_count: number
    expires: number
}

export interface ShareResponse {
    id: number
    name: string
    enable: boolean
    access_count: number
    max_access_count: number
    expires: number
    token: string
    gen: GenConfig
}

export interface ShareFilter {
    sub_id: number[]
    sub_id_exclude: boolean
    speed_up_more: number
    speed_down_more: number
    country: string[]
    country_exclude: boolean
    delay_less_than: number
    alive_status: number
    risk_less_than: number
}