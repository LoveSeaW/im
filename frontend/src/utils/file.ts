export function fileSizeFormat(size: number, pointLength: number = 2) {
    const units = ['B', 'K', 'M', 'G', 'TB'];
    let unit: string = ""
    while ((unit = units.shift() as string) && size > 1024) {
        size = size / 1024;
    }
    return (unit === 'B' ? size : size.toFixed(pointLength || 2)) + unit
}


export const base64UrlToFile = (base64Url: string, filename: string) => {
    const arr = base64Url.split(',');
    const mime = (arr[0].match(/:(.*?);/) as string[])[1];
    const bstr = atob(arr[1]);
    let n = bstr.length;
    const u8arr = new Uint8Array(n);
    while (n--) {
        u8arr[n] = bstr.charCodeAt(n);
    }
    return new File([u8arr], filename, { type: mime });
};


export function formatSecond(seconds: number) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const remainingSeconds = seconds % 60;

    const formattedHours = hours < 10 ? '0' + hours : hours;
    const formattedMinutes = minutes < 10 ? '0' + minutes : minutes;
    const formattedSeconds = remainingSeconds < 10 ? '0' + remainingSeconds : remainingSeconds;

    return `${formattedHours}:${formattedMinutes}:${formattedSeconds}`;
}