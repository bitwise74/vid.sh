type Settings = {
    defaults: {
        targetSize: number
        losslessExport: boolean
        mergeAudio: boolean
    }
    autoCopy: boolean
    disableMainPage: boolean
    richEmbeds: boolean
    preferredTheme: string
    change(key: string, val: any): void
}

type DeepPartial<T> = { [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P] }

export const DefaultSettings: Settings = {
    defaults: {
        targetSize: 20,
        losslessExport: false,
        mergeAudio: false
    },
    autoCopy: false,
    disableMainPage: false,
    richEmbeds: false,
    preferredTheme: 'system',
    change(key: string, val: any) {
        ;(this as any)[key] = val
    }
}

const validThemes = ['light', 'dark', 'system']

// Clamps settings to valid values and fills in missing values with defaults
export const clampSettings = (settings: DeepPartial<Settings>) => {
    const cs: Settings = { ...DefaultSettings, ...settings } as Settings

    // Clamp target size between 0 and 2000
    cs.defaults.targetSize = Math.min(Math.max(cs.defaults.targetSize, 0), 2000)
    cs.preferredTheme = validThemes.includes(cs.preferredTheme) ? cs.preferredTheme : DefaultSettings.preferredTheme

    return cs
}
