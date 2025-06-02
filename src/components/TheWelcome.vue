<script setup>
import Navbar from './Navbar.vue'
</script>

<template>
    <Navbar />
    <div class="card daily-task-card">
        <div class="d-flex align-items-center">
            <div class="task-text">
                <h5>Survey Kepuasan Masyarakat</h5>
                <small>Survei ini bertujuan untuk mengetahui tingkat kepuasan masyarakat terhadap layanan yang diberikan oleh Kantor Kecamatan Pulaulaut Sigam</small>
            </div>
        </div>
        <div class="progress">
            <div class="progress-bar" role="progressbar" style="width: 70%;"></div>
        </div>
    </div>
    <div class="section-title">Survey</div>
    <div class="categories-wrapper">
        <RouterLink :to="{ name: 'about' }" class="category-item text-reset text-decoration-none">
            <div class="category-icon rounded-4 fw-bold">
                <i class="bi bi-pen-fill me-3"></i>
                Lanjutkan Survey
            </div>
        </RouterLink>
    </div>
    <div class="section-title">Pelayanan Lain</div>
    <div class="more-games-wrapper">
        <div class="game-card">
            <div class="game-icon">
                <img src="https://img.icons8.com/ios-filled/100/8b5cf6/send.png" alt="Language Quiz" />
            </div>
            <div>
                <h6>Sipaktusarah</h6>
                <small>15 Surat</small>
            </div>
            <div class="card-footer">
                <div class="user-count">
                    <i class="bi bi-crown-fill"></i>
                    <span>24.7K</span>
                </div>
                <button class="play-btn">
                    <i class="bi bi-lightning-fill"></i>
                </button>
            </div>
        </div>

        <div class="game-card">
            <div class="game-icon">
                <img src="https://img.icons8.com/ios-filled/100/8b5cf6/coins.png" alt="Exam Quiz" />
            </div>
            <div>
                <h6>SiLIREK</h6>
                <small>12 Kegiatan</small>
            </div>
            <div class="card-footer">
                <div class="user-count">
                    <i class="bi bi-crown-fill"></i>
                    <span>12.5K</span>
                </div>
                <button class="play-btn">
                    <i class="bi bi-lightning-fill"></i>
                </button>
            </div>
        </div>
    </div>
</template>
