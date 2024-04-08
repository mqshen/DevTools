export class GeneralPreferences {
    theme: string
    language: string
    font: string
    fontFamily: Array<string>
    fontSize: number
    scanSize: number
    keyIconStyle: number
    useSysProxy: boolean
    useSysProxyHttp: boolean
    checkUpdate: boolean
    skipVersion: string
    allowTrack: boolean

    constructor() {
        this.theme = 'auto';
        this.language = 'auto';
        this.font = '';
        this.fontFamily = [];
        this.fontSize = 14;
        this.scanSize = 3000;
        this.keyIconStyle = 0;
        this.useSysProxy = false;
        this.useSysProxyHttp = false;
        this.checkUpdate = true;
        this.skipVersion = '';
        this.allowTrack = true;
    }
}

export class BehaviorPreferences {
    welcomed: boolean
    asideWidth: number
    windowWidth: number
    windowHeight: number
    windowMaximised: boolean

    constructor() {
        this.welcomed = false;
        this.asideWidth = 200;
        this.windowWidth = 0;
        this.windowHeight = 0;
        this.windowMaximised = false;
    }
}