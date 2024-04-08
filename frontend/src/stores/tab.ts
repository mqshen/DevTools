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

        }
    }

})


export default usePreferencesStore