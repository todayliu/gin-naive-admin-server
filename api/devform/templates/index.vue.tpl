<script setup lang="tsx">
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import type { [[.EntityName]]ListSearchParams, [[.EntityName]]Row } from './index.api'
import { Icon } from '@iconify/vue'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal } from 'naive-ui'
import { createDiscreteApi } from 'naive-ui'
import { createProSearchForm, renderProDateText } from 'pro-naive-ui'
import { computed, ref } from 'vue'
import { useProNDataTable } from '@/composables/use-pro-n-data-table'
import { Api } from './index.api'

const { message } = createDiscreteApi(['message'])

const searchForm = createProSearchForm()

const {
  search: { proSearchFormProps },
  table: { tableProps, onChange },
} = useProNDataTable(
  async ({ current, pageSize }, values) => {
    const res = await Api.list({
      pageNo: current,
      pageSize,
      ...(values as Record<string, unknown>),
    } as [[.EntityName]]ListSearchParams & { pageNo: number, pageSize: number })
    return res.data ?? { list: [], total: 0 }
  },
  { form: searchForm },
)

const searchColumns = computed<ProSearchFormColumns<[[.EntityName]]ListSearchParams>>(() => [
[[range .QueryFieldsTS -]]
  { title: '[[.Comment]]', path: '[[.JSONName]]' },
[[end]]
])

const tableColumns = computed<ProDataTableColumns<[[.EntityName]]Row>>(() => [
  { title: '序号', type: 'index', width: 70, align: 'center' },
  { title: 'ID', path: 'id', width: 80, align: 'center' },
[[range .ListColumns -]]
  { title: '[[.Title]]', path: '[[.JSONName]]', width: 160[[if .Ellipsis]], ellipsis: { tooltip: true }[[end]] },
[[end]]
[[if .HasTimeType -]]
  {
    title: '更新时间',
    width: 180,
    align: 'center',
    render: row => renderProDateText(row.updateTime),
  },
[[end]]
  {
    key: 'operation',
    title: '操作',
    width: 160,
    align: 'center',
    titleAlign: 'center',
    fixed: 'right',
    className: 'operation-column-center',
    render: (row) => {
      return (
        <n-flex justify="center" class="w-full">
          <n-button type="primary" size="small" text={true} onClick={() => openEdit(row)}>
            编辑
          </n-button>
          <n-popconfirm
            onPositiveClick={async () => {
              await Api.del(row.id)
              message.success('已删除')
              onChange({ page: 1 })
            }}
          >
            {{
              default: () => (
                <span>
                  确定删除
                  <span class="c-red-500 font-bold">{String(row.id)}</span>
                  ？
                </span>
              ),
              trigger: () => (
                <n-button type="error" size="small" text={true}>
                  删除
                </n-button>
              ),
            }}
          </n-popconfirm>
        </n-flex>
      )
    },
  },
])

const showModal = ref(false)
const formModel = ref<Partial<[[.EntityName]]Row>>({})

function openAdd() {
  formModel.value = {}
  showModal.value = true
}

function openEdit(row: [[.EntityName]]Row) {
  formModel.value = { ...row }
  showModal.value = true
}

async function submitForm() {
  try {
    const v = formModel.value
    if (v.id)
      await Api.edit(v as [[.EntityName]]Row)
    else
      await Api.add(v as Parameters<typeof Api.add>[0])

    message.success('保存成功')
    showModal.value = false
    onChange({ page: 1 })
  }
  catch {
    message.error('保存失败')
  }
}
</script>

<template>
  <n-flex
    class="h-full"
    vertical
    size="large"
  >
[[if .ShowSearch]]
    <pro-card content-class="pb-0!">
      <pro-search-form
        :form="searchForm"
        :columns="searchColumns"
        v-bind="proSearchFormProps"
        label-width="auto"
      />
    </pro-card>
[[end]]
    <pro-data-table
      :title="'[[.ListTitle]]'"
      row-key="id"
      :columns="tableColumns"
      v-bind="tableProps"
      :scroll-x="960"
      :flex-height="false"
    >
      <template #toolbar>
        <n-flex>
          <n-button type="primary" @click="openAdd">
            <template #icon>
              <n-icon>
                <Icon icon="ant-design:plus-outlined" />
              </n-icon>
            </template>
            新增
          </n-button>
        </n-flex>
      </template>
    </pro-data-table>

    <NModal
      v-model:show="showModal"
      preset="card"
      :title="formModel.id != null && formModel.id !== undefined ? '编辑' : '新增'"
      class="max-w-560px"
    >
      <NForm :model="formModel" label-placement="left" label-width="auto">
        <NFormItem v-if="formModel.id != null && formModel.id !== undefined" label="ID">
          <span>{{ formModel.id }}</span>
        </NFormItem>
[[range .FormFieldsVue -]]
        <NFormItem label="[[.Label]]">
[[.ControlTpl]]
        </NFormItem>
[[end]]
        <n-flex justify="end">
          <NButton @click="showModal = false">
            取消
          </NButton>
          <NButton type="primary" @click="submitForm">
            保存
          </NButton>
        </n-flex>
      </NForm>
    </NModal>
  </n-flex>
</template>

<style lang="scss" scoped>
:deep(.operation-column-center) {
  text-align: center;

  .n-flex {
    justify-content: center;
    width: 100%;
  }
}
</style>
