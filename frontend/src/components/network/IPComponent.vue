<script setup lang="ts" async>
import { NotificationReactive } from "naive-ui";
import { setupDiscreteApi, MyNotificationApi } from '@/utils/discrete'
import { GetIP } from "wailsjs/go/services/ipService.js";
import { ref, onMounted } from "vue";

const ips = ref([]);
const loading = ref(true);

const getIp = () => {
    GetIP().then((resp) => {
        if (resp.success) {
            ips.value = resp.data.ips;
        }
        loading.value = false;
    });
};
const iconClass = (idx: number) => {
    return ["bi", `bi-${idx + 1}-circle-fill`];
};
const notification = ref<MyNotificationApi | null>(null)
const copyIP = async (ip: string) => {
    if (!notification.value) {
        notification.value = await setupDiscreteApi()
    }
    navigator.clipboard.writeText(ip).then(function () {
        notification.value?.success('Success', {
            title: 'Copy',
            meta: "Copy Success"
        })
    }, function (err) {
        console.error('Async: Could not copy text: ', err);
    });
}
onMounted(() => {
    getIp();
});
</script>

<template>
    <n-spin :show="loading" :theme-overrides="{ opacitySpinning: 0 }">
        <n-grid x-gap="12" :cols="3">
            <n-gi v-for="(ip, i) in ips" :key="i" class="ip-card">
                <n-card>
                    <template #header>
                        <div class="ip-card-header" style="font-weight: bold">
                            <span>
                                <i :class="iconClass(i)"></i>
                                <span class="title-text">
                                    IP 来源: {{ ip.source }}
                                </span>
                            </span>
                            <n-button quaternary>
                                <i class="bi bi-arrow-clockwise"></i>
                                <!-- <n-icon><cash-icon /></n-icon> -->
                            </n-button>
                        </div>
                        <div class="p-3 ip-link">
                            <span>
                                <i class="bi bi-pc-display-horizontal"></i>
                            </span>
                            <span class="title-text">
                                {{ ip.ip }}

                                <n-button quaternary @click="copyIP(ip.ip)">
                                    <i class="bi bi-clipboard-plus" role="button"></i>
                                    <!-- <n-icon><cash-icon /></n-icon> -->
                                </n-button>
                            </span>
                        </div>
                    </template>
                    <div class="ip-card-body">
                        <ul class="list-group list-group-flush">
                            <li>
                                <span>
                                    <i class="bi bi-geo-alt-fill"></i>
                                    地区 :
                                </span>
                                <span v-if="ip.location">
                                    {{ ip.location.country }}
                                </span>
                            </li>
                            <li class="">
                                <span>
                                    <i class="bi bi-houses"></i> 省份 :
                                </span v-if="ip.location">
                                <span> {{ ip.location.province }}</span>
                            </li>
                            <li class="">
                                <span>
                                    <i class="bi bi-sign-turn-right"></i> 城市 :
                                </span>
                                <span v-if="ip.location"> {{ ip.location.city }}</span>
                            </li>
                            <li class="">
                                <span>
                                    <i class="bi bi-ethernet"></i> 网络 :
                                </span>
                                <span v-if="ip.location"> {{ ip.location.isp }}</span>
                            </li>
                            <!-- <li class="">
              <span >
                <i class="bi bi-reception-4"></i> 类型 :
              </span>
              <span >无线网络 </span>
            </li> -->
                            <!-- <li class="j">
              <span >
                <i class="bi bi-shield-fill-check"></i> 代理 :
              </span>
              <span >不是代理或 VPN </span>
            </li> -->
                            <li class="j border-0">
                                <span>
                                    <i class="bi bi-buildings"></i> ASN :
                                </span>

                                <span> {{ ip.location.asn }}</span>
                                <span class="col-9">
                                    <i class="bi bi-info-circle" role="button">
                                    </i>
                                </span>
                            </li>
                        </ul>
                    </div>
                </n-card>
            </n-gi>
        </n-grid>
    </n-spin>
</template>

<style lang="scss">
.ip-card-header {
    display: flex;
    justify-content: space-between;
    flex-direction: row;
}

.ip-link {
    padding: 20px 0 0 0;

    &::before {
        content: "";
        position: absolute;
        top: 45px;
        left: 34px;
        transform: translate(-50%);
        height: 25px;
        width: 2px;
        z-index: 1;
        border-left: 2px dashed rgb(33, 37, 41);
    }
}

.n-card>.n-card-header {
    background-color: var(--n-action-color);

    .title-text {
        margin-left: 5px;
    }
}

.ip-card {
    margin-bottom: 20px;
}

.ip-card-body {
    ul {
        padding: 0px;
    }

    li {
        display: flex;
        border-bottom: 1px solid rgb(222, 226, 230);
        padding: 12px 4px;
        margin: 0px;
    }
}
</style>
