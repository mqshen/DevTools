import { defineStore } from 'pinia'
import { TabItem } from '@/objects/tabItem'


const usePreferencesStore = defineStore('tab', {
    state: () => ({
        tabList: new Array<TabItem>(),
        activatedIndex: 0, 
    }), 
    getters: {
        /**
         * get current tab list item
         * @returns {TabItem[]}
         */
        tabs(): Array<TabItem> {
            return this.tabList
        },
        currentTabName(): string {
            return this.tabs[this.activatedIndex].name
        },
    },
    actions: {
        /**
         *
         * @param idx
         * @param {boolean} [switchNav]
         * @param {string} [subTab]
         * @private
         */
        _setActivatedIndex(idx: number) {
            this.activatedIndex = idx
        },
        upsertTab(name: string, compontent: Object) {
            let tabIndex = this.tabList.findIndex((item: TabItem) => item.name === name)
            if (tabIndex === -1) {
                const tabItem = new TabItem(name, compontent)
                this.tabList.push(tabItem)
                tabIndex = this.tabList.length - 1
            }
            this._setActivatedIndex(tabIndex)

        },
        closeTab(name: string) {
            const idx = this.tabs.findIndex((item: TabItem) => item.name === name)
            if (idx !== -1) {
                this.removeTab(idx)
            }
        },
        /**
         *
         * @param {number} tabIndex
         * @returns {*|null}
         */
        removeTab(tabIndex: number) {
            const len = this.tabs.length

            if (tabIndex < 0 || tabIndex >= len) {
                return null
            }
            const removed = this.tabList.splice(tabIndex, 1)

            // update select index if removed index equal current selected
            this.activatedIndex -= 1
            if (this.activatedIndex < 0) {
                if (this.tabList.length > 0) {
                    this._setActivatedIndex(0)
                } else {
                    this._setActivatedIndex(-1)
                }
            } else {
                this._setActivatedIndex(this.activatedIndex)
            }

            return removed.length > 0 ? removed[0] : null
        },
    }

})


export default usePreferencesStore