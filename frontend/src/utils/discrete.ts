import usePreferencesStore from "stores/preferences";
import {
  createDiscreteApi,
  darkTheme,
  NotificationApi,
  NotificationOptions,
  NotificationReactive,
} from "naive-ui";
import { darkThemeOverrides, themeOverrides } from "@/utils/theme";
import { i18nGlobal } from "@/utils/i18n";
import { computed } from "vue";

export interface MyNotificationApi {
  show: (option: NotificationOptions) => NotificationReactive;
  error: (content: string, option: any) => NotificationReactive;
  info: (content: string, option: any) => NotificationReactive;
  success: (content: string, option: any) => NotificationReactive;
  warning: (content: string, option: any) => NotificationReactive;
}

function setupNotification(notification: NotificationApi) {
  return {
    /**
     * @param {NotificationOption} option
     * @return {NotificationReactive}
     */
    show(option: NotificationOptions) {
      return notification.create(option);
    },
    error: (content: string, option: any = {}) => {
      option.content = content;
      option.title = option.title || i18nGlobal.t("common.error");
      return notification.error(option);
    },
    info: (content: string, option: any = {}) => {
      option.content = content;
      return notification.info(option);
    },
    success: (content: string, option: any = {}) => {
      option.content = content;
      option.title = option.title || i18nGlobal.t("common.success");
      return notification.success(option);
    },
    warning: (content: string, option: any = {}) => {
      option.content = content;
      option.title = option.title || i18nGlobal.t("common.warning");
      return notification.warning(option);
    },
  };
}

/**
 *
 * @param {DialogApiInjection} dialog
 * @return {*}
 */
function setupDialog(dialog) {
  return {
    /**
     * @param {DialogOptions} option
     * @return {DialogReactive}
     */
    show(option) {
      option.closable = option.closable === true;
      option.autoFocus = option.autoFocus === true;
      option.transformOrigin = "center";
      return dialog.create(option);
    },
    warning: (content, onConfirm) => {
      return dialog.warning({
        title: i18nGlobal.t("common.warning"),
        content: content,
        closable: false,
        autoFocus: false,
        transformOrigin: "center",
        positiveText: i18nGlobal.t("common.confirm"),
        negativeText: i18nGlobal.t("common.cancel"),
        onPositiveClick: () => {
          onConfirm && onConfirm();
        },
      });
    },
  };
}

/**
 * setup discrete api and bind global component (like dialog, message, alert) to window
 * @return {Promise<void>}
 */
export async function setupDiscreteApi() {
  const prefStore = usePreferencesStore();
  const configProviderProps = computed(() => ({
    theme: prefStore.isDark ? darkTheme : undefined,
    themeOverrides,
  }));
  const { message, dialog, notification } = createDiscreteApi(
    ["message", "notification", "dialog"],
    {
      configProviderProps,
      messageProviderProps: {
        placement: "bottom",
        keepAliveOnHover: true,
        containerStyle: {
          marginBottom: "38px",
        },
        themeOverrides: prefStore.isDark
          ? darkThemeOverrides.Message
          : themeOverrides.Message,
      },
      notificationProviderProps: {
        max: 5,
        placement: "bottom-right",
        keepAliveOnHover: true,
        containerStyle: {
          marginBottom: "38px",
        },
      },
    },
  );

  const $notification = setupNotification(notification);
  window.$notification = $notification;
  window.$dialog = setupDialog(dialog);
  return $notification;
}
