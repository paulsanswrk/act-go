<script setup lang="ts">

import {ref} from "vue";
import axios from "axios";

const username = ref()
const pwd = ref()


async function submit() {
  const {data:{success}} = await axios.post('http://localhost:8080/login', {username: username.value, pwd: pwd.value})

  console.log({success})
}

</script>

<template>
  <div
      class="px-6 py-20 md:px-12 lg:px-20 flex items-center justify-center bg-[linear-gradient(-225deg,var(--p-primary-500),var(--p-primary-700)_48%,var(--p-primary-800))] dark:bg-[linear-gradient(-225deg,var(--p-primary-400),var(--p-primary-600)_48%,var(--p-primary-800))]"
  >
    <form method="POST" action="http://localhost/login" @submit.prevent="submit"
          class="p-12 shadow text-center lg:w-[30rem] backdrop-blur-md rounded-xl bg-[rgba(255,255,255,0.1)]">
      <div class="text-4xl font-medium mb-12 text-primary-contrast">ACT</div>
      <InputText
          type="text"
          v-model.trim="username"
          class="!appearance-none placeholder:!text-primary-contrast/40 !border-0 !p-4 !w-full !outline-0 !text-xl !block !mb-6 !bg-white/10 !text-primary-contrast/70 !rounded-full"
          placeholder="Email"
      />
      <InputText
          type="password"
          v-model.trim="pwd"
          class="!appearance-none placeholder:!text-primary-contrast/40 !border-0 !p-4 !w-full !outline-0 !text-xl !mb-6 !bg-white/10 !text-primary-contrast/70 !rounded-full"
          placeholder="Password"
      />
      <button
          type="submit"
          class="max-w-40 w-full rounded-full appearance-none border-0 p-4 outline-0 text-xl mb-6 font-medium bg-white/30 hover:bg-white/40 active:bg-white/20 text-primary-contrast/80 cursor-pointer transition-colors duration-150"
      >
        Sign In
      </button>
      <a class="cursor-pointer font-medium block text-center text-primary-contrast">Forgot Password?</a>
    </form>
  </div>
</template>

<style>

</style>
