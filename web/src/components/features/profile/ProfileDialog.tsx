"use client"

import { useState, useEffect, useCallback } from "react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"

import {
  Dialog,
  DialogContent,
} from "@/src/components/ui/dialog"
import { Button } from "@/src/components/ui/button"
import { Input } from "@/src/components/ui/input"
import { Label } from "@/src/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/src/components/ui/card"
import { Badge } from "@/src/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/src/components/ui/table"
import { ProfileLayout } from "./ProfileLayout"
import { InlineLoading } from "@/src/components/ui/loading"
import { api } from "@/src/lib/api/client"
import { useAuth } from "@/src/components/providers"
import type { SessionInfo } from "@/src/types"

interface ProfileDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

interface FormData {
  username: string
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

export function ProfileDialog({ open, onOpenChange }: ProfileDialogProps) {
  const { user, logout, updateUser } = useAuth()
  const [activeTab, setActiveTab] = useState("profile")
  const [sessions, setSessions] = useState<SessionInfo[]>([])
  const [isLoadingSessions, setIsLoadingSessions] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const form = useForm<FormData>({
    mode: 'onChange',
    defaultValues: {
      username: user?.username || "",
      oldPassword: "",
      newPassword: "",
      confirmPassword: ""
    }
  })


  useEffect(() => {
    if (open && user) {
      form.reset({
        username: user.username,
        oldPassword: "",
        newPassword: "",
        confirmPassword: ""
      })
      loadSessions()
    }
  }, [open, user])

  const loadSessions = useCallback(async () => {
    setIsLoadingSessions(true)
    try {
      const response = await api.getSessions()
      setSessions(response.sessions)
    } catch (error) {
      console.error('Failed to load sessions:', error)
      toast.error("加载会话列表失败")
    } finally {
      setIsLoadingSessions(false)
    }
  }, [])

  const handleUpdateUsername = useCallback(async (data: FormData) => {
    if (!data.username.trim()) {
      toast.error("用户名不能为空")
      return
    }

    if (data.username === user?.username) {
      toast.error("新用户名不能与当前用户名相同")
      return
    }

    setIsSubmitting(true)
    try {
      await api.updateUsername({ username: data.username })
      updateUser({ ...user!, username: data.username })
      toast.success("用户名修改成功")
      onOpenChange(false)
    } catch (error: any) {
      toast.error(error.message || "用户名修改失败")
    } finally {
      setIsSubmitting(false)
    }
  }, [user, updateUser, onOpenChange])

  const handleChangePassword = useCallback(async (data: FormData) => {
    if (!data.oldPassword || !data.newPassword) {
      toast.error("请填写完整密码信息")
      return
    }

    if (data.newPassword !== data.confirmPassword) {
      toast.error("两次输入的新密码不一致")
      return
    }

    if (data.newPassword.length < 6) {
      toast.error("新密码长度至少为6位")
      return
    }

    setIsSubmitting(true)
    try {
      await api.changePassword({
        username: user!.username,
        old_password: data.oldPassword,
        new_password: data.newPassword
      })
      toast.success("密码修改成功，请重新登录")
      onOpenChange(false)
      // 密码修改成功后调用登出
      setTimeout(() => {
        logout()
      }, 1000)
    } catch (error: any) {
      toast.error(error.message || "密码修改失败")
    } finally {
      setIsSubmitting(false)
    }
  }, [user, logout, onOpenChange])

  const handleDeleteSession = useCallback(async (sessionId: number) => {
    try {
      await api.deleteSession(sessionId)
      toast.success("会话已删除")
      loadSessions()
    } catch (error: any) {
      toast.error(error.message || "删除会话失败")
    }
  }, [loadSessions])

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault()

    const data = form.getValues()

    if (activeTab === "profile") {
      await handleUpdateUsername(data)
    } else if (activeTab === "password") {
      await handleChangePassword(data)
    }
  }, [activeTab, form, handleUpdateUsername, handleChangePassword])

  const getBrowserInfo = (userAgent: string) => {
    const ua = userAgent.toLowerCase()

    // Edge first (since it includes Chrome)
    if (ua.includes('edg/')) return 'Edge'
    if (ua.includes('edge/')) return 'Edge'

    // Opera (includes Chrome)
    if (ua.includes('opr/')) return 'Opera'
    if (ua.includes('opera')) return 'Opera'

    // Chrome (check before Safari, but after Edge/Opera)
    if (ua.includes('chrome/') && !ua.includes('edg/') && !ua.includes('opr/')) return 'Chrome'

    // Firefox
    if (ua.includes('firefox')) return 'Firefox'

    // Safari (must be checked after Chrome)
    if (ua.includes('safari') && !ua.includes('chrome')) return 'Safari'

    return userAgent.split(' ')[0]
  }

  const getOSInfo = (userAgent: string) => {
    const ua = userAgent.toLowerCase()

    // Check for iOS devices first (before macOS)
    if (ua.includes('iphone')) return 'iPhone'
    if (ua.includes('ipad')) return 'iPad'
    if (ua.includes('ipod')) return 'iPod'

    if (ua.includes('windows nt 10.0')) return 'Windows 10'
    if (ua.includes('windows nt 11.0')) return 'Windows 11'
    if (ua.includes('mac os x')) return 'macOS'
    if (ua.includes('linux')) return 'Linux'
    if (ua.includes('android')) return 'Android'

    return 'Unknown OS'
  }

  const renderContent = () => {
    switch (activeTab) {
      case "profile":
        return (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="username">用户名</Label>
              <Input
                id="username"
                {...form.register("username", { required: true })}
                placeholder="请输入用户名"
              />
              {form.formState.errors.username && (
                <p className="text-sm text-destructive">用户名不能为空</p>
              )}
            </div>
          </div>
        )

      case "password":
        return (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="oldPassword">当前密码</Label>
              <Input
                id="oldPassword"
                type="password"
                {...form.register("oldPassword", { required: true })}
                placeholder="请输入当前密码"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="newPassword">新密码</Label>
              <Input
                id="newPassword"
                type="password"
                {...form.register("newPassword", {
                  required: true,
                  minLength: 6
                })}
                placeholder="请输入新密码（至少6位）"
              />
              {form.formState.errors.newPassword && (
                <p className="text-sm text-destructive">密码长度至少为6位</p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">确认新密码</Label>
              <Input
                id="confirmPassword"
                type="password"
                {...form.register("confirmPassword", {
                  required: true,
                  validate: (value) => value === form.getValues("newPassword") || "两次输入的密码不一致"
                })}
                placeholder="请再次输入新密码"
              />
              {form.formState.errors.confirmPassword && (
                <p className="text-sm text-destructive">{form.formState.errors.confirmPassword.message}</p>
              )}
            </div>
          </div>
        )

      case "sessions":
        return (
          <div className="space-y-4">
            <div className="lg:hidden space-y-3">
              {isLoadingSessions ? (
                <div className="flex items-center justify-center py-8">
                  <InlineLoading message="加载会话中..." size="sm" />
                </div>
              ) : sessions.length === 0 ? (
                <p className="text-center text-muted-foreground py-8">暂无其他活跃会话</p>
              ) : (
                sessions.map((session) => (
                  <Card key={session.id}>
                    <CardContent className="p-4">
                      <div className="flex items-start justify-between">
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-2">
                            <span className="font-medium truncate">
                              {getBrowserInfo(session.user_agent)}
                            </span>
                            <Badge variant={session.is_active ? "default" : "secondary"} className="text-xs">
                              {session.is_active ? "活跃" : "已失效"}
                            </Badge>
                          </div>
                          <div className="space-y-1 text-sm text-muted-foreground">
                            <div className="flex items-center gap-2">
                              <span className="text-xs">系统:</span>
                              <span className="text-xs">{getOSInfo(session.user_agent)}</span>
                            </div>
                            <div className="flex items-center gap-2">
                              <span className="text-xs">IP:</span>
                              <span className="font-mono text-xs">{session.client_ip}</span>
                            </div>
                            <div className="flex items-center gap-2">
                              <span className="text-xs">创建:</span>
                              <span className="text-xs">{new Date(session.created_at).toLocaleDateString()}</span>
                            </div>
                            <div className="flex items-center gap-2">
                              <span className="text-xs">最后活跃:</span>
                              <span className="text-xs">{new Date(session.last_access_at).toLocaleString()}</span>
                            </div>
                          </div>
                        </div>
                        <div className="ml-3 flex-shrink-0">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleDeleteSession(session.id)}
                            disabled={!session.is_active}
                            className="h-8 px-3 text-xs"
                          >
                            注销
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </div>

            <Card className="hidden lg:block">
              <CardHeader>
                <CardTitle className="text-base">活跃会话</CardTitle>
                <CardDescription>
                  管理您在所有设备上的登录会话
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoadingSessions ? (
                  <div className="flex items-center justify-center py-8">
                    <InlineLoading message="加载会话中..." size="sm" />
                  </div>
                ) : sessions.length === 0 ? (
                  <p className="text-center text-muted-foreground py-8">暂无其他活跃会话</p>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>设备</TableHead>
                        <TableHead>IP地址</TableHead>
                        <TableHead>状态</TableHead>
                        <TableHead>最后活跃</TableHead>
                        <TableHead className="w-[100px]">操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {sessions.map((session) => (
                        <TableRow key={session.id}>
                          <TableCell className="font-medium">
                            <div className="flex flex-col">
                              <span>{getBrowserInfo(session.user_agent)}</span>
                              <span className="text-xs text-muted-foreground">
                                {getOSInfo(session.user_agent)}
                              </span>
                            </div>
                          </TableCell>
                          <TableCell>{session.client_ip}</TableCell>
                          <TableCell>
                            <Badge variant={session.is_active ? "default" : "secondary"}>
                              {session.is_active ? "活跃" : "已失效"}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            {new Date(session.last_access_at).toLocaleString()}
                          </TableCell>
                          <TableCell>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleDeleteSession(session.id)}
                              disabled={!session.is_active}
                            >
                              注销
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </div>
        )

      default:
        return null
    }
  }

  const renderActions = (isMobile?: boolean) => {
    if (activeTab === "sessions") {
      return (
        <Button
          type="button"
          variant="outline"
          onClick={() => onOpenChange(false)}
          className={isMobile ? "flex-1 h-9 sm:h-10 text-sm" : "h-10"}
        >
          关闭
        </Button>
      )
    }

    const isDirty = form.formState.isDirty
    const canSubmit = activeTab === "profile"
      ? isDirty && form.getValues("username").trim() !== user?.username
      : isDirty && form.getValues("newPassword") === form.getValues("confirmPassword") && form.getValues("newPassword").length >= 6

    return (
      <>
        <Button
          type="button"
          variant="outline"
          onClick={() => onOpenChange(false)}
          className={isMobile ? "flex-1 h-9 sm:h-10 text-sm" : "h-10"}
        >
          取消
        </Button>
        <Button
          type="submit"
          disabled={!canSubmit || isSubmitting}
          className={isMobile ? "flex-1 h-9 sm:h-10 text-sm" : "h-10"}
        >
          {isSubmitting ? <InlineLoading size="sm" /> : "保存"}
        </Button>
      </>
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="overflow-hidden p-0 md:max-h-[600px] md:max-w-[800px] lg:max-w-[900px] max-h-[90vh] h-full md:h-auto w-[95vw] sm:w-[90vw] md:w-full">
        <ProfileLayout
          activeTab={activeTab}
          onTabChange={setActiveTab}
          onSubmit={handleSubmit}
          renderActions={renderActions}
        >
          {renderContent()}
        </ProfileLayout>
      </DialogContent>
    </Dialog>
  )
}