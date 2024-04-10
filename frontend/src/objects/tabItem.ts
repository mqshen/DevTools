/**
 * tab item
 */
export class TabItem {
    name: string
    compontent: string
    icon: string

    /**
     *
     * @param {string} name connection name
     */
    constructor(name: string, compontent: string, icon: string) {
        this.name = name
        this.compontent = compontent
        this.icon = icon
    }
}
