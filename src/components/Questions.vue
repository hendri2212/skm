<template>
    <div class="question-wrapper mx-auto" style="max-width: 414px;">
        <!-- ===== FORM DATA PENGGUNA ===== -->
        <div v-if="!quizStarted" id="user-form">
            <div class="d-flex justify-content-between align-items-center" style="margin-bottom: 1.5rem;">
                <button class="btn btn-link p-0 bg-white rounded-circle" style="width: 48px; height: 48px;"
                    @click="backToHome()">
                    <i class="bi bi-arrow-left fs-4 text-secondary"></i>
                </button>
                <h6 class="mb-0">Data Pengguna</h6>
                <button class="btn p-0"></button>
            </div>
            <div class="card question-card mb-4 rounded-4 shadow-sm">
                <div class="card-body">
                    <!-- Nama Lengkap -->
                    <div class="form-floating mb-3">
                        <input v-model="user.name" type="text" class="form-control" id="input-name"
                            placeholder="Nama Lengkap" required />
                        <label for="input-name">Nama Lengkap</label>
                    </div>

                    <div class="row">
                        <!-- Tempat Lahir -->
                        <div class="col-md-6 mb-3">
                            <div class="form-floating">
                                <input v-model="user.birthPlace" type="text" class="form-control" id="input-birth-place"
                                    placeholder="Tempat Lahir" required />
                                <label for="input-birth-place">Tempat Lahir</label>
                            </div>
                        </div>

                        <!-- Tanggal Lahir -->
                        <div class="col-md-6 mb-3">
                            <div class="form-floating">
                                <input v-model="user.birthDate" type="date" class="form-control" id="input-birth-date"
                                    placeholder="Tanggal Lahir" required />
                                <label for="input-birth-date">Tanggal Lahir</label>
                            </div>
                        </div>
                    </div>

                    <!-- Jenis Kelamin (Radio Inputs) -->
                    <div class="mb-3">
                        <label class="form-label d-block">Jenis Kelamin</label>
                        <div class="form-check form-check-inline">
                            <input class="form-check-input" type="radio" name="gender" id="gender-male"
                                value="Laki-laki" v-model="user.gender" required />
                            <label class="form-check-label" for="gender-male">Laki-laki</label>
                        </div>
                        <div class="form-check form-check-inline">
                            <input class="form-check-input" type="radio" name="gender" id="gender-female"
                                value="Perempuan" v-model="user.gender" required />
                            <label class="form-check-label" for="gender-female">Perempuan</label>
                        </div>
                    </div>

                    <!-- Pendidikan Terakhir -->
                    <div class="form-floating mb-3">
                        <select v-model="user.education" class="form-select" id="select-education" required>
                            <option value="" disabled>— Pilih Pendidikan Terakhir —</option>
                            <option value="SD Kebawah">SD Kebawah</option>
                            <option value="SLTP">SLTP</option>
                            <option value="SLTA">SLTA</option>
                            <option value="S1-D3-D4">S1-D3-D4</option>
                            <option value="S-2 ke atas">S-2 ke atas</option>
                        </select>
                        <label for="select-education">Pendidikan Terakhir</label>
                    </div>

                    <!-- Pekerjaan Utama -->
                    <div class="form-floating mb-3">
                        <select v-model="user.occupation" class="form-select" id="select-occupation" required>
                            <option value="" disabled>— Pilih Pekerjaan Utama —</option>
                            <option value="PNS/TNI/Polri">PNS/TNI/Polri</option>
                            <option value="Pegawai Swasta">Pegawai Swasta</option>
                            <option value="Wiraswastawan/Usahawan">Wiraswastawan/Usahawan</option>
                            <option value="Pelajar/Mahasiswa">Pelajar/Mahasiswa</option>
                            <option value="Lainnya">Lainnya</option>
                        </select>
                        <label for="select-occupation">Pekerjaan Utama</label>
                    </div>

                    <!-- Tombol Mulai Quiz -->
                    <button class="btn btn-warning w-100 mt-5" :disabled="!canStartQuiz" @click="startQuiz">
                        Mulai Quiz
                    </button>
                </div>
            </div>
        </div>

        <!-- ===== LAYAR SOAL ===== -->
        <div v-else-if="!showResult" class="quiz-screen">
            <!-- Header: Nomor Soal -->
            <div class="d-flex justify-content-between align-items-center" style="margin-bottom: 1.5rem;">
                <button class="btn btn-link p-0 bg-white rounded-circle" @click="prevQuestion" style="width: 48px; height: 48px;">
                    <i class="bi bi-arrow-left fs-4 text-secondary"></i>
                </button>
                <h6 class="mb-0">
                    Soal {{ currentIndex + 1 }} / {{ questions.length }}
                </h6>
                <button class="btn p-0"></button>
            </div>

            <!-- Kartu Soal Utama -->
            <div class="card question-card mb-4 shadow-sm rounded-4">
                <div class="card-body">
                    <!-- Teks Soal -->
                    <p class="question-text p-5 rounded-4 text-center">{{ currentQuestion.question }}</p>

                    <!-- Progress Bar Manual -->
                    <div class="mb-3">
                        <div class="time-bar-container">
                            <div class="time-bar-fill"
                                :style="{ width: ((currentIndex + 1) / questions.length * 100) + '%' }"></div>
                        </div>
                    </div>

                    <!-- Pilihan Jawaban (A, B, C, D) -->
                    <div class="answer-list">
                        <div class="answer-option" v-for="(label, key) in currentQuestion.options" :key="key" :class="{
                            selected: selectedOption === key
                        }" @click="selectOption(key)">
                            <div class="d-flex align-items-center w-100">
                                <!-- Huruf Pilihan -->
                                <div class="option-label me-3">{{ key.toUpperCase() }}</div>
                                <!-- Teks Pilihan -->
                                <div class="option-text">{{ label }}</div>
                            </div>
                        </div>
                    </div>

                    <button class="btn btn-warning w-100 mt-5" @click="goToNext" :disabled="!selectedOption">
                        {{ isLastQuestion ? 'Lihat Hasil' : 'Soal Berikutnya' }}
                    </button>
                </div>
            </div>

        </div>

        <!-- ===== HALAMAN HASIL AKHIR ===== -->
        <div v-else class="result-screen text-center">
            <!-- Header: Nomor Soal -->
            <div class="d-flex justify-content-between align-items-center" style="margin-bottom: 1.5rem;">
                <button class="btn btn-link p-0 bg-white rounded-circle" @click="backToHome"
                    :disabled="currentIndex === 0" style="width: 48px; height: 48px;">
                    <i class="bi bi-house fs-4 text-secondary"></i>
                </button>
            </div>
            <div class="card question-card mb-4 shadow-sm rounded-4">
                <div class="card-body">
                    <h3 class="mb-3">Quiz Selesai!</h3>
                    <p class="mb-2">
                        Halo <strong>{{ user.name }}</strong>, Anda telah menyelesaikan quiz ini!
                    </p>
                    <h4 class="text-primary">Terima kasih atas partisipasi Anda</h4>
                    <button class="btn btn-success mt-5 w-100" @click="restartQuiz">
                        Ulangi Quiz
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Questions',
    data() {
        return {
            // ===== STATE UNTUK FORM PENGGUNA =====
            user: {
                name: '',
                birthPlace: '',
                birthDate: '',
                gender: '',
                education: '',
                occupation: ''
            },

            // Apakah quiz sudah dimulai?
            quizStarted: false,

            // Daftar soal; setiap objek memuat:
            // { question: 'Teks soal', options: { a: '...', b: '...', c: '...', d: '...' } }
            questions: [
                {
                    question: 'Siapa penemu telepon pertama kali?',
                    options: {
                        a: 'Alexander Graham Bell',
                        b: 'Thomas Edison',
                        c: 'Nikola Tesla',
                        d: 'Guglielmo Marconi',
                    },
                },
                {
                    question: 'Apa ibukota Indonesia?',
                    options: {
                        a: 'Jakarta',
                        b: 'Bandung',
                        c: 'Surabaya',
                        d: 'Yogyakarta',
                    },
                },
                {
                    question: 'Berapa hasil dari 7 × 8?',
                    options: {
                        a: '45',
                        b: '56',
                        c: '63',
                        d: '49',
                    },
                },
                // Tambahkan soal sebanyak yang dibutuhkan…
            ],

            // Index soal yang sedang ditampilkan (0-based)
            currentIndex: 0,

            // Opsi yang dipilih user untuk soal sekarang (key: 'a'|'b'|'c'|'d')
            selectedOption: '',

            // Apakah sudah mencapai akhir dan tampilkan hasil?
            showResult: false,

            // Array untuk menyimpan jawaban user
            userAnswers: [],
        };
    },
    computed: {
        // Soal yang sedang aktif
        currentQuestion() {
            return this.questions[this.currentIndex];
        },
        // Apakah ini soal terakhir?
        isLastQuestion() {
            return this.currentIndex === this.questions.length - 1;
        },
        // Tombol Mulai aktif jika nama & email terisi
        canStartQuiz() {
            // return this.user.name.trim() !== '' && this.user.email.trim() !== '';
            return (
                this.user.name.trim() &&
                this.user.birthPlace.trim() &&
                this.user.birthDate &&
                this.user.gender &&
                this.user.education &&
                this.user.occupation
            );
        },
    },
    methods: {
        // Kembali ke halaman home
        backToHome() {
            this.$router.push({ name: 'home' });
        },

        // Mulai Quiz: inisialisasi state
        startQuiz() {
            this.quizStarted = true;
            this.userAnswers = [];
        },

        // User memilih opsi
        selectOption(key) {
            this.selectedOption = key;
        },

        // Pindah ke soal berikutnya / tampilkan hasil jika terakhir
        goToNext() {
            // Simpan jawaban user
            this.userAnswers[this.currentIndex] = this.selectedOption;

            this.selectedOption = '';
            if (!this.isLastQuestion) {
                this.currentIndex++;
            } else {
                // Quiz selesai
                this.showResult = true;
            }
        },

        // Pindah ke soal sebelumnya (opsional)
        prevQuestion() {
            if (this.currentIndex === 0) {
                // Perubahan: kembali ke form data pengguna
                this.quizStarted = false;
            } else {
                this.currentIndex--;
                // Ambil jawaban yang sudah dipilih sebelumnya jika ada
                this.selectedOption = this.userAnswers[this.currentIndex] || '';
            }
        },

        // Ulangi quiz dari awal
        restartQuiz() {
            this.quizStarted = false;
            this.showResult = false;
            this.currentIndex = 0;
            this.selectedOption = '';
            this.userAnswers = [];
        },
    },
};
</script>

<style scoped>
/* ===================== STYLING UNTUK Questions.vue ===================== */

/* .question-text: styling judul soal */
.question-text {
    font-size: 1.1rem;
    font-weight: 600;
    color: #ffffff;
    margin-bottom: 1rem;
    background: linear-gradient(135deg, #8b5cf6 0%, #d8b4fe 100%);
}

/* ===== PROGRESS BAR ===== */
.time-bar-container {
    height: 8px;
    background-color: #e5e7eb;
    border-radius: 4px;
    overflow: hidden;
}

.time-bar-fill {
    background-color: #fbbf24;
    height: 100%;
    transition: width 0.3s ease;
}

/* ===== LIST JAWABAN ===== */
.answer-list {
    display: flex;
    flex-direction: column;
}

.answer-option {
    background-color: #f3f4f6;
    border-radius: 0.75rem;
    padding: 0.75rem 1rem;
    margin-bottom: 0.75rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border: 2px solid transparent;
    transition: background-color 0.2s, border-color 0.2s;
}

.answer-option:hover {
    background-color: #e5e7eb;
}

/* Pilihan yang dipilih */
.answer-option.selected {
    background-color: #cffafe;
    border-color: #0ea5e9;
}

/* Huruf opsi (A, B, C, D) */
.option-label {
    font-weight: 600;
    font-size: 1rem;
    color: #1f2937;
    width: 24px;
    text-align: center;
}

/* Teks opsi */
.option-text {
    font-size: 0.95rem;
    color: #374151;
    flex: 1;
}
</style>