import { Controller, Control, UseFormWatch, UseFormReset } from 'react-hook-form'
import { useState, useEffect, useCallback } from 'react'
import { Input } from '@/src/components/ui/input'
import { Label } from '@/src/components/ui/label'
import { Switch } from '@/src/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/src/components/ui/select'
import { Calendar22 } from '../share-date-pick'
import { SUB_RULES } from '../../constants/sub-rules'
import { SUBSCRIPTION_TARGETS } from '../../constants'
import { isCustomConfig } from '../../utils'
import type { ShareRequest, KeyValue } from '@/src/types'

interface ConfigSectionProps {
    control: Control<ShareRequest>
    watch: UseFormWatch<ShareRequest>
    reset: UseFormReset<ShareRequest>
}

export function ConfigSection({ control, watch, reset }: ConfigSectionProps) {
    const currentSubConverter = watch('gen.sub_converter')
    const currentConfig = currentSubConverter?.config || ''

    const defaultRuleValue = SUB_RULES[0]?.value || ''

    const getCurrentSelection = useCallback(() => {
        if (!currentConfig) return defaultRuleValue
        if (isCustomConfig(currentConfig, SUB_RULES)) return 'custom'
        return currentConfig
    }, [currentConfig, defaultRuleValue])

    const [currentSelection, setCurrentSelection] = useState(() => getCurrentSelection())

    useEffect(() => {
        if (!currentConfig && defaultRuleValue) {
            const formData = watch()
            const currentValues = formData.gen?.sub_converter || {}
            reset({
                ...formData,
                gen: {
                    ...formData.gen,
                    sub_converter: {
                        ...currentValues,
                        config: defaultRuleValue
                    }
                }
            })
        }
    }, [])

    const handleRuleChange = useCallback((value: string) => {
        setCurrentSelection(value)

        const formData = watch()
        const currentValues = formData.gen?.sub_converter || {}

        const newConfig = value === 'custom' ? '' : value

        reset({
            ...formData,
            gen: {
                ...formData.gen,
                sub_converter: {
                    ...currentValues,
                    config: newConfig
                }
            }
        })
    }, [watch, reset])

    return (
        <div className="space-y-4">
            <div>
                <Label htmlFor="template" className="mb-2 block">
                    订阅模板
                </Label>
                <Controller
                    name="gen.sub_converter.target"
                    control={control}
                    render={({ field }) => (
                        <Select onValueChange={field.onChange} value={field.value ?? 'auto'}>
                            <SelectTrigger className="w-full">
                                <SelectValue placeholder="选择订阅模板" />
                            </SelectTrigger>
                            <SelectContent>
                                {SUBSCRIPTION_TARGETS.map((target) => (
                                    <SelectItem key={target.value} value={target.value}>
                                        {target.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    )}
                />
            </div>

            <div>
                <Label htmlFor="rename" className="mb-2 block">
                    重命名模板
                </Label>
                <Controller
                    name="gen.rename"
                    control={control}
                    render={({ field }) => (
                        <Input
                            {...field}
                            value={field.value || ''}
                            id="rename"
                        />
                    )}
                />
            </div>

            <div>
                <Label htmlFor="sub_converter_config" className="mb-2 block">
                    规则链接
                </Label>
                <Select onValueChange={handleRuleChange} value={currentSelection}>
                    <SelectTrigger className="w-full" id="sub_converter_config">
                        <SelectValue placeholder="选择规则链接" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="custom">自定义</SelectItem>
                        {SUB_RULES.length > 0 ? (
                            SUB_RULES.map((item: KeyValue) => (
                                <SelectItem key={item.value} value={item.value}>
                                    {item.key}
                                </SelectItem>
                            ))
                        ) : (
                            <SelectItem value="" disabled>
                                暂无可用规则
                            </SelectItem>
                        )}
                    </SelectContent>
                </Select>

                {currentSelection === 'custom' && (
                    <div className="mt-2">
                        <Controller
                            name="gen.sub_converter.config"
                            control={control}
                            render={({ field }) => (
                                <Input
                                    {...field}
                                    value={field.value || ''}
                                    id="sub_converter_config_custom"
                                    placeholder="规则链接"
                                />
                            )}
                        />
                    </div>
                )}
            </div>

            <div className="flex items-center justify-between">
                <Label htmlFor="proxy">代理</Label>
                <Controller
                    name="gen.proxy"
                    control={control}
                    render={({ field }) => (
                        <Switch
                            id="proxy"
                            checked={!!field.value}
                            onCheckedChange={field.onChange}
                        />
                    )}
                />
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div>
                    <Label htmlFor="max_access_count" className="mb-2 block">
                        最大访问次数
                    </Label>
                    <Controller
                        name="max_access_count"
                        control={control}
                        render={({ field }) => (
                            <Input
                                {...field}
                                value={field.value || ''}
                                id="max_access_count"
                                type="number"
                                placeholder="0"
                                min="0"
                                onChange={(e) => field.onChange(parseInt(e.target.value || '0'))}
                            />
                        )}
                    />
                </div>

                <div>
                    <Controller
                        name="expires"
                        control={control}
                        render={({ field }) => (
                            <Calendar22
                                value={field.value ?? 0}
                                onChange={(ts: number) => field.onChange(ts || 0)}
                            />
                        )}
                    />
                </div>
            </div>
        </div>
    )
}