<template>
    <button 
        :class="[
            'inline-flex flex-none w-auto items-center justify-center rounded-full transition-colors',
            variant === 'icon' ? 'p-2 hover:bg-gray-100' : 'px-4 py-2 hover:bg-gray-600'
        ]"
        v-bind="$attrs"
    >
        <slot></slot>
    </button>
</template>

<script>
export default {
    name: 'BaseButton',
    props: {
        variant: {
            type: String,
            default: 'default',
            validator: value => ['default', 'icon'].includes(value)
        }
    }
}
</script>
