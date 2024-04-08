/**
 * tab item
 */
export class TabItem {
    name: string
    compontent: Object

    /**
     *
     * @param {string} name connection name
     */
    constructor(name: string, compontent: string) {
        this.name = name
        this.compontent = compontent
    }
}
