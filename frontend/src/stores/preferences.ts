import { defineStore } from 'pinia'
import { cloneDeep, get, set, isEmpty, join, map, pick } from 'lodash'
import { GeneralPreferences, BehaviorPreferences } from '@/objects/general'
import { useOsTheme, enUS, zhCN } from 'naive-ui'
import { i18nGlobal } from '@/utils/i18n'
import { GetPreferences, SetPreferences } from 'wailsjs/go/services/preferencesService.js'

const osTheme = useOsTheme()

const usePreferencesStore = defineStore('preferences', {
    state: () => ({
        behavior: new BehaviorPreferences(),
        general: new GeneralPreferences(),
        lastPref: {},
    }), 
    getters: {
        themeLocale() {
            const lang = this.currentLanguage
            switch (lang) {
                case 'zh':
                    return zhCN
                default:
                    return enUS
            }
        },
        /**
         * get current language setting
         * @return {string}
         */
        currentLanguage(): string {
            let lang = this.general.language
            if (lang === 'auto') {
                const systemLang = navigator.language //|| navigator.userLanguage
                lang = systemLang.split('-')[0]
            }
            return lang || 'en'
        },
        isDark(): boolean {
            const th = this.general.theme
            if (th !== 'auto') {
                return th === 'dark'
            }
            return osTheme.value === 'dark'
        },
        /**
         * current font selection
         * @returns {{fontSize: string, fontFamily?: string}}
         */
        generalFont(): string {
            const fontStyle: any = {
                fontSize: this.general.fontSize + 'px',
            }
            if (!isEmpty(this.general.fontFamily)) {
                fontStyle['fontFamily'] = join(
                    map(this.general.fontFamily, (f) => `"${f}"`),
                    ',',
                )
            }
            return fontStyle
        },
    },
    actions: {
        _applyPreferences(data: any) {
            for (const key in data) {
                set(this, key, data[key])
            }
        },
        /**
         * load preferences from local
         * @returns {Promise<void>}
         */
        async loadPreferences() {
            const { success, data } = await GetPreferences()
            if (success) {
                this.lastPref = cloneDeep(data)
                this._applyPreferences(data)
                // default value
                const showLineNum = get(data, 'editor.showLineNum')
                if (showLineNum === undefined) {
                    set(data, 'editor.showLineNum', true)
                }
                const showFolding = get(data, 'editor.showFolding')
                if (showFolding === undefined) {
                    set(data, 'editor.showFolding', true)
                }
                const dropText = get(data, 'editor.dropText')
                if (dropText === undefined) {
                    set(data, 'editor.dropText', true)
                }
                const links = get(data, 'editor.links')
                if (links === undefined) {
                    set(data, 'editor.links', true)
                }
                i18nGlobal.locale.value = this.currentLanguage
            }
        },
        /**
         * save preferences to local
         * @returns {Promise<boolean>}
         */
        async savePreferences() {
            const pf = pick(this, ['behavior', 'general', 'editor', 'cli', 'decoder'])
            const { success, msg } = await SetPreferences(pf)
            return success === true
        },
    }

})


export default usePreferencesStore