<script setup>
import Photo from '../components/Photo.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import ErrorMsg from '../components/ErrorMsg.vue'
</script>

<script>
export default {
    components: {
        LoadingSpinner,
        ErrorMsg,
    },
    data: function () {
        return {
            errormsg: null,
			loading: false,
            username: null,
        }
    },
    methods: {
        
    },
    mounted() {
        if (this.$profile.identifier == null) {
            this.$router.push({ path: '/' });
        } else {
            console.log(this.username);
        }  
    }
}
</script>

<template>
    <div class="d-flex justify-content-center align-items-start min-vh-100 w-100" style="background-color: #383838">
        <div class="d-flex flex-column justify-content-centeralign-items-start min-vh-100 w-75 text-white" style="background-color: #2e2e2e;">		
			<div class="cerca pt-70 text-white fs-1">
                {{ this.username }}
            </div>
		</div>
    </div>
</template>

<style scoped>


</style>
