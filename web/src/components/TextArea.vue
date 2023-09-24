<script setup>
defineProps({
  label: String,
})
</script>
<template>
  <div class="form-control">
    <label class="label">
        <span class="label-text">{{label}}</span>
    </label>
    <textarea class="textarea textarea-bordered w-full h-full border-2" readonly></textarea>
  </div>
</template>
