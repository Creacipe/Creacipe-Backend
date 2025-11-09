-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Nov 09, 2025 at 07:24 AM
-- Server version: 8.0.30
-- PHP Version: 8.1.10

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `creacipepintar_db`
--

-- --------------------------------------------------------

--
-- Table structure for table `categories`
--

CREATE TABLE `categories` (
  `category_id` bigint UNSIGNED NOT NULL,
  `category_name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `categories`
--

INSERT INTO `categories` (`category_id`, `category_name`) VALUES
(1, 'Bahan Utama'),
(2, 'Metode Masak'),
(3, 'Jenis Hidangan'),
(4, 'Rasa');

-- --------------------------------------------------------

--
-- Table structure for table `comments`
--

CREATE TABLE `comments` (
  `comment_id` bigint UNSIGNED NOT NULL,
  `menu_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `comment_text` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `log_activity`
--

CREATE TABLE `log_activity` (
  `activity_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `action` varchar(100) NOT NULL COMMENT 'e.g., CREATE_MENU, LOGIN, DELETE_USER',
  `target_id` bigint UNSIGNED DEFAULT NULL,
  `target_type` varchar(100) DEFAULT NULL COMMENT 'e.g., menus, users, comments',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `log_activity`
--

INSERT INTO `log_activity` (`activity_id`, `user_id`, `action`, `target_id`, `target_type`, `created_at`) VALUES
(1, 1, 'SETUP_FIRST_ADMIN', 1, 'users', '2025-10-23 13:55:32'),
(2, 2, 'USER_REGISTER', 2, 'users', '2025-10-23 13:56:42'),
(3, 3, 'USER_REGISTER', 3, 'users', '2025-10-23 13:57:06'),
(4, 4, 'USER_REGISTER', 4, 'users', '2025-10-23 13:57:25'),
(5, 1, 'ADMIN_CREATE_USER', 5, 'users', '2025-10-23 13:58:08'),
(6, 1, 'USER_LOGIN', 1, 'users', '2025-10-23 13:58:58'),
(7, 1, 'ADMIN_DEACTIVATE_USER', 5, 'users', '2025-10-23 14:00:00'),
(8, 1, 'ADMIN_ACTIVATE_USER', 5, 'users', '2025-10-23 14:00:18'),
(9, 1, 'ADMIN_UPDATE_ROLE', 4, 'users', '2025-10-23 14:02:16'),
(10, 1, 'ADMIN_DEACTIVATE_USER', 5, 'users', '2025-10-23 14:16:15'),
(11, 5, 'USER_LOGIN', 5, 'users', '2025-10-23 14:16:55'),
(12, 5, 'USER_LOGIN', 5, 'users', '2025-10-23 14:17:27'),
(13, 1, 'ADMIN_ACTIVATE_USER', 5, 'users', '2025-10-23 14:30:02'),
(14, 5, 'USER_LOGIN', 5, 'users', '2025-10-23 14:30:20'),
(15, 2, 'USER_LOGIN', 2, 'users', '2025-10-23 14:35:09'),
(16, 2, 'USER_LOGIN', 2, 'users', '2025-10-23 14:35:12'),
(17, 2, 'CREATE_MENU', 1, 'menus', '2025-10-23 15:09:44'),
(18, 2, 'CREATE_MENU', 2, 'menus', '2025-10-23 15:17:49'),
(19, 2, 'CREATE_MENU', 3, 'menus', '2025-10-23 15:18:12'),
(20, 2, 'CREATE_MENU', 4, 'menus', '2025-10-23 15:18:52'),
(21, 2, 'CREATE_MENU', 5, 'menus', '2025-10-23 15:19:15'),
(22, 2, 'CREATE_MENU', 6, 'menus', '2025-10-23 15:19:42'),
(23, 2, 'CREATE_MENU', 7, 'menus', '2025-10-23 15:20:07'),
(24, 2, 'CREATE_MENU', 8, 'menus', '2025-10-23 15:20:29'),
(25, 2, 'CREATE_MENU', 9, 'menus', '2025-10-23 15:21:07'),
(26, 3, 'USER_LOGIN', 3, 'users', '2025-10-23 15:23:03'),
(27, 3, 'USER_LOGIN', 3, 'users', '2025-10-23 15:23:06'),
(28, 3, 'CREATE_MENU', 10, 'menus', '2025-10-23 15:24:27'),
(29, 3, 'CREATE_MENU', 11, 'menus', '2025-10-23 15:25:16'),
(30, 3, 'CREATE_MENU', 12, 'menus', '2025-10-23 15:25:31'),
(31, 3, 'CREATE_MENU', 13, 'menus', '2025-10-23 15:25:45'),
(32, 3, 'CREATE_MENU', 14, 'menus', '2025-10-23 15:26:04'),
(33, 3, 'CREATE_MENU', 15, 'menus', '2025-10-23 15:26:19'),
(34, 3, 'CREATE_MENU', 16, 'menus', '2025-10-23 15:27:04'),
(35, 3, 'CREATE_MENU', 17, 'menus', '2025-10-23 15:27:23'),
(36, 4, 'USER_LOGIN', 4, 'users', '2025-10-23 15:32:01'),
(37, 4, 'CREATE_MENU', 18, 'menus', '2025-10-23 15:33:32'),
(38, 4, 'CREATE_MENU', 19, 'menus', '2025-10-23 15:34:04'),
(39, 4, 'CREATE_MENU', 20, 'menus', '2025-10-23 15:34:15'),
(40, 4, 'CREATE_MENU', 21, 'menus', '2025-10-23 15:34:27'),
(41, 4, 'CREATE_MENU', 22, 'menus', '2025-10-23 15:34:45'),
(42, 4, 'CREATE_MENU', 23, 'menus', '2025-10-23 15:34:58'),
(43, 4, 'CREATE_MENU', 24, 'menus', '2025-10-23 15:35:11'),
(44, 4, 'CREATE_MENU', 25, 'menus', '2025-10-23 15:35:58'),
(45, 4, 'CREATE_MENU', 26, 'menus', '2025-10-23 16:23:16'),
(46, 1, 'USER_LOGIN', 1, 'users', '2025-10-26 05:11:08'),
(47, 1, 'CREATE_CATEGORY', 5, 'categories', '2025-10-26 05:20:36'),
(48, 1, 'UPDATE_CATEGORY', 5, 'categories', '2025-10-26 05:24:33'),
(49, 1, 'DELETE_CATEGORY', 5, 'categories', '2025-10-26 05:25:47'),
(50, 1, 'CREATE_TAG', 55, 'tags', '2025-10-26 05:27:43'),
(51, 1, 'UPDATE_TAG', 55, 'tags', '2025-10-26 05:29:47'),
(52, 1, 'DELETE_TAG', 55, 'tags', '2025-10-26 05:31:19'),
(53, 1, 'APPROVE_MENU', 26, 'menus', '2025-10-26 06:02:33'),
(54, 1, 'REJECT_MENU', 25, 'menus', '2025-10-26 06:04:56'),
(55, 1, 'USER_LOGIN', 1, 'users', '2025-10-26 06:17:15'),
(56, 6, 'USER_REGISTER', 6, 'users', '2025-10-26 06:25:26'),
(57, 6, 'USER_LOGIN', 6, 'users', '2025-10-26 06:25:42'),
(58, 6, 'CREATE_MENU', 27, 'menus', '2025-10-26 06:29:49'),
(59, 6, 'CREATE_MENU', 28, 'menus', '2025-10-26 06:30:44'),
(60, 6, 'UPDATE_MENU', 28, 'menus', '2025-10-26 06:34:29'),
(61, 1, 'USER_LOGIN', 1, 'users', '2025-10-26 06:35:13'),
(62, 1, 'UPDATE_MENU', 28, 'menus', '2025-10-26 06:35:32'),
(63, 6, 'USER_LOGIN', 6, 'users', '2025-10-26 06:37:27'),
(64, 6, 'DELETE_MENU', 28, 'menus', '2025-10-26 06:38:26'),
(65, 1, 'USER_LOGIN', 1, 'users', '2025-10-26 06:38:59'),
(66, 1, 'DELETE_MENU', 27, 'menus', '2025-10-26 06:39:23'),
(67, 6, 'USER_LOGIN', 6, 'users', '2025-10-26 06:43:41'),
(68, 6, 'BOOKMARK_MENU', 1, 'menus', '2025-10-26 06:45:39'),
(69, 6, 'UNBOOKMARK_MENU', 1, 'menus', '2025-10-26 06:46:38'),
(70, 6, 'LIKE_MENU', 1, 'menus', '2025-10-26 06:47:56'),
(71, 6, 'DISLIKE_MENU', 1, 'menus', '2025-10-26 06:49:01'),
(72, 6, 'REMOVE_VOTE', 1, 'menus', '2025-10-26 06:50:38'),
(73, 6, 'LIKE_MENU', 1, 'menus', '2025-10-26 06:53:28'),
(74, 6, 'REMOVE_VOTE', 1, 'menus', '2025-10-26 07:03:33'),
(75, 6, 'LIKE_MENU', 1, 'menus', '2025-10-26 07:03:53'),
(76, 6, 'LIKE_MENU', 2, 'menus', '2025-10-26 07:04:47'),
(77, 6, 'BOOKMARK_MENU', 1, 'menus', '2025-10-26 07:05:06'),
(78, 6, 'BOOKMARK_MENU', 2, 'menus', '2025-10-26 07:07:50'),
(79, 6, 'BOOKMARK_MENU', 3, 'menus', '2025-10-26 07:07:54'),
(80, 2, 'USER_LOGIN', 2, 'users', '2025-10-26 07:08:18'),
(81, 2, 'BOOKMARK_MENU', 3, 'menus', '2025-10-26 07:08:39'),
(82, 2, 'BOOKMARK_MENU', 5, 'menus', '2025-10-26 07:08:43'),
(83, 2, 'BOOKMARK_MENU', 1, 'menus', '2025-10-26 07:08:48'),
(84, 2, 'BOOKMARK_MENU', 7, 'menus', '2025-10-26 07:08:52'),
(85, 6, 'REMOVE_VOTE', 1, 'menus', '2025-10-26 07:09:06'),
(86, 6, 'REMOVE_VOTE', 2, 'menus', '2025-10-26 07:09:10'),
(87, 6, 'LIKE_MENU', 3, 'menus', '2025-10-26 07:09:15'),
(88, 6, 'LIKE_MENU', 4, 'menus', '2025-10-26 07:09:21'),
(89, 3, 'USER_LOGIN', 3, 'users', '2025-10-26 07:09:47'),
(90, 3, 'BOOKMARK_MENU', 1, 'menus', '2025-10-26 07:10:03'),
(91, 3, 'BOOKMARK_MENU', 2, 'menus', '2025-10-26 07:10:07'),
(92, 3, 'LIKE_MENU', 1, 'menus', '2025-10-26 07:10:22'),
(93, 3, 'LIKE_MENU', 2, 'menus', '2025-10-26 07:10:26'),
(94, 3, 'LIKE_MENU', 3, 'menus', '2025-10-26 07:10:30'),
(95, 4, 'USER_LOGIN', 4, 'users', '2025-10-26 07:10:49'),
(96, 4, 'BOOKMARK_MENU', 2, 'menus', '2025-10-26 07:11:01'),
(97, 4, 'BOOKMARK_MENU', 1, 'menus', '2025-10-26 07:11:05'),
(98, 4, 'BOOKMARK_MENU', 3, 'menus', '2025-10-26 07:11:10'),
(99, 4, 'LIKE_MENU', 3, 'menus', '2025-10-26 07:11:18'),
(100, 4, 'LIKE_MENU', 2, 'menus', '2025-10-26 07:11:22'),
(101, 4, 'LIKE_MENU', 1, 'menus', '2025-10-26 07:11:26'),
(102, 1, 'USER_LOGIN', 1, 'users', '2025-10-26 07:20:48'),
(103, 1, 'APPROVE_MENU', 1, 'menus', '2025-10-26 07:21:06'),
(104, 1, 'APPROVE_MENU', 2, 'menus', '2025-10-26 07:21:11'),
(105, 1, 'APPROVE_MENU', 3, 'menus', '2025-10-26 07:21:15'),
(106, 1, 'APPROVE_MENU', 4, 'menus', '2025-10-26 07:21:19'),
(107, 1, 'APPROVE_MENU', 5, 'menus', '2025-10-26 07:21:25'),
(108, 1, 'APPROVE_MENU', 6, 'menus', '2025-10-26 07:21:28'),
(109, 1, 'APPROVE_MENU', 7, 'menus', '2025-10-26 07:21:34'),
(110, 1, 'APPROVE_MENU', 8, 'menus', '2025-10-26 07:21:38'),
(111, 1, 'APPROVE_MENU', 9, 'menus', '2025-10-26 07:21:47'),
(112, 1, 'APPROVE_MENU', 10, 'menus', '2025-10-26 07:21:51'),
(113, 1, 'APPROVE_MENU', 11, 'menus', '2025-10-26 07:21:55'),
(114, 1, 'APPROVE_MENU', 12, 'menus', '2025-10-26 07:21:59'),
(115, 1, 'APPROVE_MENU', 13, 'menus', '2025-10-26 07:22:03'),
(116, 1, 'APPROVE_MENU', 14, 'menus', '2025-10-26 07:22:07'),
(117, 1, 'APPROVE_MENU', 15, 'menus', '2025-10-26 07:22:11'),
(118, 1, 'APPROVE_MENU', 16, 'menus', '2025-10-26 07:22:16'),
(119, 1, 'APPROVE_MENU', 17, 'menus', '2025-10-26 07:22:22'),
(120, 1, 'APPROVE_MENU', 18, 'menus', '2025-10-26 07:22:27'),
(121, 1, 'APPROVE_MENU', 19, 'menus', '2025-10-26 07:22:31'),
(122, 1, 'APPROVE_MENU', 20, 'menus', '2025-10-26 07:22:37'),
(123, 7, 'USER_REGISTER', 7, 'users', '2025-10-26 14:11:46'),
(124, 7, 'USER_LOGIN', 7, 'users', '2025-10-26 14:12:09'),
(125, 7, 'LIKE_MENU', 1, 'menus', '2025-10-26 14:14:10'),
(126, 7, 'BOOKMARK_MENU', 1, 'menus', '2025-10-26 14:14:30'),
(127, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 09:51:48'),
(128, 2, 'LIKE_MENU', 5, 'menus', '2025-11-08 09:54:39'),
(129, 2, 'DISLIKE_MENU', 5, 'menus', '2025-11-08 09:54:48'),
(130, 2, 'BOOKMARK_MENU', 5, 'menus', '2025-11-08 09:54:50'),
(131, 2, 'LIKE_MENU', 1, 'menus', '2025-11-08 10:10:45'),
(132, 2, 'REMOVE_VOTE', 1, 'menus', '2025-11-08 10:10:46'),
(133, 2, 'LIKE_MENU', 1, 'menus', '2025-11-08 10:10:46'),
(134, 2, 'REMOVE_VOTE', 1, 'menus', '2025-11-08 10:10:47'),
(135, 2, 'LIKE_MENU', 1, 'menus', '2025-11-08 10:10:48'),
(136, 2, 'REMOVE_VOTE', 1, 'menus', '2025-11-08 10:10:48'),
(137, 2, 'LIKE_MENU', 1, 'menus', '2025-11-08 10:10:49'),
(138, 2, 'REMOVE_VOTE', 1, 'menus', '2025-11-08 10:10:49'),
(139, 2, 'LIKE_MENU', 1, 'menus', '2025-11-08 10:10:49'),
(140, 2, 'REMOVE_VOTE', 1, 'menus', '2025-11-08 10:10:50'),
(141, 2, 'BOOKMARK_MENU', 1, 'menus', '2025-11-08 10:10:50'),
(142, 2, 'UNBOOKMARK_MENU', 1, 'menus', '2025-11-08 10:10:52'),
(143, 2, 'BOOKMARK_MENU', 1, 'menus', '2025-11-08 10:10:52'),
(144, 2, 'UNBOOKMARK_MENU', 1, 'menus', '2025-11-08 10:10:53'),
(145, 2, 'BOOKMARK_MENU', 11, 'menus', '2025-11-08 10:10:56'),
(146, 2, 'LIKE_MENU', 11, 'menus', '2025-11-08 10:10:58'),
(147, 2, 'BOOKMARK_MENU', 10, 'menus', '2025-11-08 10:11:07'),
(148, 2, 'LIKE_MENU', 10, 'menus', '2025-11-08 10:11:08'),
(149, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 10:13:50'),
(150, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 11:38:36'),
(151, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 11:38:59'),
(152, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 11:39:33'),
(153, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 11:40:20'),
(154, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 11:40:59'),
(155, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 11:41:05'),
(156, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 12:03:36'),
(157, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 12:05:24'),
(158, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 12:15:19'),
(159, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 12:15:38'),
(160, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 12:31:31'),
(161, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 12:47:39'),
(162, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 12:48:15'),
(163, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 12:54:31'),
(164, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 13:03:28'),
(165, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 13:18:24'),
(166, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 13:24:26'),
(167, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 13:46:32'),
(168, 2, 'USER_LOGIN', 2, 'users', '2025-11-08 13:47:52'),
(169, 2, 'CREATE_MENU', 29, 'menus', '2025-11-09 06:02:17'),
(170, 2, 'USER_LOGIN', 2, 'users', '2025-11-09 06:49:39');

-- --------------------------------------------------------

--
-- Table structure for table `menus`
--

CREATE TABLE `menus` (
  `menu_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text,
  `ingredients` json DEFAULT NULL,
  `instructions` json DEFAULT NULL,
  `image_url` varchar(255) DEFAULT NULL,
  `status` enum('pending','approved','rejected') NOT NULL DEFAULT 'pending',
  `rejection_reason` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `menus`
--

INSERT INTO `menus` (`menu_id`, `user_id`, `title`, `description`, `ingredients`, `instructions`, `image_url`, `status`, `rejection_reason`, `created_at`, `updated_at`) VALUES
(1, 2, 'Ayam Bumbu Bawang', 'Resep ayam goreng sederhana dengan marinasi bumbu bawang yang gurih dan meresap.', '[\"1 dada ayam fillet, cuci bersih, potong kecil2\", \"1 buah jeruk nipis\", \"Bumbu marinate:\", \"5 siung bawang putih\", \"1 sdt lada hitam\", \"secukupnya Gula dan garam\", \"Minyam goreng\"]', '[\"Kucuri ayam dengan air jeruk nipis, kemudian marinate dengan bumbu halus selama 3 jam di kulkas\", \"Panaskan minyak, goreng ayam sampai matang dan kecoklatan\", \"Sajikan hangat dengan nasi putih.\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:09:44', '2025-10-26 07:21:06'),
(2, 2, 'Ayam Woku Manado', 'Resep Ayam Woku Manado yang lezat dan mudah dibuat.', '[\"1 Ekor Ayam Kampung (potong 12)--2 Buah Jeruk Nipis--2 Sdm Garam--3 Ruas Kunyit--7 Bawang Merah--7 Bawang Putih--10 Cabe Merah--10 Cabe Rawit Merah (sesuai selera)--3 Butir Kemiri--2 Batang Sereh--2 Lembar Daun Salam--2 Ikat Daun Kemangi--Penyedap Rasa--1 1/2 Gelas Air--\"]', '[\"Cuci bersih ayam dan tiriskan. Lalu peras jeruk nipis (kalo gak ada jeruk nipis bisa pake cuka) dan beri garam. Aduk hingga merata dan diamkan selama 5 menit\", \"biar ayam gak bau amis.--Goreng ayam tersebut setengah matang\", \"lalu tiriskan--Haluskan bumbu menggunakan blender. Bawang merah\", \"bawang putih\", \"cabe merah\", \"cabe rawit\", \"kemiri dan kunyit. Oh iya kasih minyak sedikit yaa biar bisa di blender. Untuk sereh nya di geprek aja terus di buat simpul.--Setelah bumbu di haluskan barulah di tumis. Jangan lupa sereh dan daun salamnya juga ikut di tumis. Di tumis sampai berubah warna ya 👌--Masukan ayam yang sudah di goreng setengah matang ke dalam bumbu yang sudah di tumis\", \"dan diamkan 5 menit dulu. Biar bumbu meresap. Lalu tuangkan 1 1/2 Gelas air. Lalu tambahkan penyedap rasa (saya 3 Sdt\", \"tapi sesuai selera ya) koreksi rasa dan Biar kan sampai mendidih--Setelah masakan mendidih\", \"lalu masukan daun kemangi yang sudah di potong potong. Masak lagi sekitar 10 menit. And taraaaaaaaaaaaaaa..... jadi deh Ayam Woku Manadonya.--Oh iyaa kalo mau di tambahkan potongan tomat merah juga bisa ko. Sesuai selera aja yaa buibuuuu 👌👌👌--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:17:49', '2025-10-26 07:21:11'),
(3, 2, 'Gurame Saus Padang', 'Resep Gurame Saus Padang yang lezat dan mudah dibuat.', '[\"Bahan utama:--1 ekor gurame--Bumbu untuk saus:--4 siung bawang putih (cincang halus)--3 siung bawang merah (cincang halus)--15 bh cabai merah (giling) banyaknya cabai sesuai selera ya😊--7 buah cabai rawit iris tipis (sesuai selera)--1/2 bawang bombang iris--Saus tiram--Saus tomat--Garam--Gula--Lada--1 bh wortel iris tipis memanjang (opsional)--1 buah tomat potong dadu (opsional)--1 batang irisan daun bawang--1 ruas jahe ukuran kecil iris tipis--Tepung maizena--250 ml air--Bumbu untuk menggoreng ikan :--Tepung beras--Tepung terigu--Garam--Jeruk nipis/ lemon--Minyak sayur--\"]', '[\"Cuci bersih ikan gurame yang akan dimasak. Setelah itu lumuri ikan dengan air perasan jeruk nipis/lemon diamkan beberapa menit untuk menghilangkan bau amisnya--Campur tepung terigu\", \"tepung beras dan garam menjadi satu\", \"kemudian lumuri ikan dengan tepung tersebut lalu goreng hingga matang lalu tiriskan. Pastikan minyak yang digunakan sudah panas ya biar ikannya tidak lengket😊.--Tumis bawang merah bawang putih hingga harum\", \"kemudian masukan cabai yang sudah giling\", \"aduk rata--Tambahkan bawang bombay\", \"saus tiram saus tomat dan jahe\", \"aduk sebentar lalu tambahkan 250ml air--Kemudian tambahkan garam dan gula sesuai selera.--Masukkan irisan cabai rawit dan daun bawang dan tomat. Aduk rata. Kemudian tambahkan larutan maizena agar saus mengental--Jika tidak ada tepung maizena. Diamkan beberapa menit hingga kandungan air menyusut dan saus terlihat mengental--Tes rasa. Jika sudah pas. Sajikan dan gurame saus padang siap disantap!😊--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:18:12', '2025-10-26 07:21:15'),
(4, 2, 'Sate Kambing', 'Resep Sate Kambing yang lezat dan mudah dibuat.', '[\"Bahan-bahan :--500 gr Daging kambing--Daun pepaya--Bumbu yang dihaluskan :--4 siung bamer--4 siung baput--2 cm jahe--1 sdt ketumbar--1/2 sdt lada bubuk--1/2 sdt Garam--1 sdm air asam jawa--Bahan sambel kecap :--Kecap manis sesuai selera (aku pake bango)--Cabe rawit--Bamer--Tomat--Jeruk limau--\"]', '[\"1. Cuci bersih daging kambing\", \"potong\\\" kotak\", \"bungkus didalam lapisan daun pepaya (30menit)--2.setelah 30 menit\", \"keluarkan daging dr daun pepaya\", \"campurkan dengan bumbu yang sudah dihaluskan\", \"diamkan 15 menit--3. Susun daging pada tusuk sate kurlen 4-5 potong per tusuk--4. Bakar (bisa pakai arang / teflon)--5. Bolak balik sate hingga benar\\\" matang--6. Sajikan dengan sambal kecap--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:18:52', '2025-10-26 07:21:19'),
(5, 2, 'Beef Teriyaki', 'Resep Beef Teriyaki yang lezat dan mudah dibuat.', '[\"250 gr daging sapi--1 siung bawang bombai--5 siung bawang putih--1 sachet saus teriyaki (sy pk yg saor*, bsa diganti saus tiram)--1 sdm kecap manis--secukupnya garam--secukupnya lada--secukupnya gula--secukupnya penyedap rasa--\"]', '[\"Potong kecil-kecil memanjang daging sapi, lalu cuci bersih--Tambahkan garam dan lada diamkan selama kurleb 15 menit--Iris bawang bombai dan bawang putih--Tumis bawang bombai dan bawang putih, setelah harum masukkan daging--Tumis daging sebentar lalu tambahkan saus teriyaki, kecap manis, garam, lada, gula dan penyedap rasa, beri air sedikit--Tutup wajannya agar panasnya merata--Tunggu sampai daging matang, lalu garnish sesuai selera--Beef teriyaki siap untuk disantap--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:19:15', '2025-10-26 07:21:25'),
(6, 2, 'Martabak tahu pedas kulit lumpia', 'Resep Martabak tahu pedas kulit lumpia yang lezat dan mudah dibuat.', '[\"12 lembar kulit lumpia--4 buah tahu putih--2 batang daun bawang (iris halus)--1 batang daun seledri (iris halus)--3 btr telur (sy pakai 2 aja krn agak besar)--1 sdm terigu + 2 sdm air (sbgai lem kulit lumpia)--Sckpnya minyak utk menggoreng--Bumbu halus :--3 siung bawang putih--1/2 bks lada/merica bubuk--Sckpnya garam gula dan kaldu bubuk--3-4 cabe merah kriting iris serong--2 buah wortel agak kecil parut kasar--\"]', '[\"Haluskan tahu putih campur semua bahan kecuali kulit lumpia (icip rasa)--Ambil 1 lbr kulit lumpia beri 1 sdm adonan tahu ditengahnya lalu lipat (seperti amplop) dgn rapi, lalu beri lem disetiap sisisnya spy tdk terbuka pd saat digoreng, lakukan hingga habis--Panaskan minyak diatas wajan dgn api sedang cenderung kecil, goreng martabak sampai kecoklatan angkat dan tiriskan (abaikan wajannya yg sdh mulai jelek ya 😁)--Martabak tahu pedas kulit lumpia siap dihidangkan--Selamat mencoba\", \"😘😘--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:19:42', '2025-10-26 07:21:28'),
(7, 2, 'Orak arik telur buncis', 'Resep Orak arik telur buncis yang lezat dan mudah dibuat.', '[\"1 buah telur--10 buncis--3 bawang merah--2 bawang putih--6 cabe rawit--1 sdm kecap manis--Penyedap (garam gula atau penyedap lainnya)--\"]', '[\"Haluskan bawang merah\", \"bawang putih dan cabai--Iris buncis sesuai selera--Tumis bumbu halus hingga harum--Masukan buncis hingga agak layu--Masukan penyedap--Masukan telur\", \"lalu campur dengan buncis di orak arik hingga rata.. tunggu sebentar lalu di orak arik lagi.. supaya matangnya sempurna 👌(biar ga amis)--Masukan kecap manis--Masukan air sedikit biarkan menyusut dan meresap--Koreksi rasa--Jadi deh.. selamat menikmati 😇--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:20:07', '2025-10-26 07:21:34'),
(8, 2, 'Orek tempe manis pedas super simple', 'Resep Orek tempe manis pedas super simple yang lezat dan mudah dibuat.', '[\"Tempe (saya beli 3000 saja)--3 buah cabe gendot--3 buah cabe keriting--3 buah cabe rawit merah--2 siung bawang merah--1 siung bawang putih--Kecap--Saus tiram--Saus tomat--Kaldu bubuk--Garam--Gula--\"]', '[\"Potong2 tempe sesuai selera lalu goreng sebentar.--Tumis duo bawang dan semua cabe (cabe bisa sesuai selera ya pedesnya suka segimana) sampai harum. Masukkan tempe\", \"kecap\", \"saus tiram\", \"saus tomat\", \"kaldu bubuk\", \"gulgar\", \"dan tambahkan air satu gelas belimbing. Aduk2 dan tutup wajan.--Saat air sudah menyusut dan kuah kecap mengental. Cek rasa. Klo sudah oke. Matikan api. Sajikan. (Klo aku suka airnya sedikit sekali 😋). Selamat mencoba !--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:20:29', '2025-10-26 07:21:38'),
(9, 2, 'Lumpia udang kulit tahu\'', 'Resep Lumpia udang kulit tahu  yang lezat dan mudah dibuat.', '[\"50 gram ayam potong kotak kecil--200 gram udang besar potong kecil--1 lembar daun bawang iris halus--1/2 sendok teh garam halus/sesuai selera--1/2 sendok teh gula--1/4 sendok teh merica bubuk--1 sendok teh kecap ikan--1 sendok teh saus raja rasa--1 sendok teh saus tiram--1 sendok teh minyak wijen--1 sendok makan tepung sagu--1 sendok teh tepung terigu--1 butir telur ambil bagian putihnya saja, kocok lepas--2 lembar kembang tahu kering / kulit tahu direndam sebentar--\"]', '[\"Campur ayam & udang dengan semua bumbu & daun bawang\", \"aduk rata. Sesudah rata masukkan tepung sagu\", \"terigu & putih telur. Aduk rata.--Isikan ke lembaran kulit kembang tahu & lipat seperti melipat lumpia\", \"sambil dipadatkan agar rapi & bentuk agak pipih jangan terlalu bulat. Olesi putih telur diujung lipatan agar merekat rapat.--Kukus sekitar 15 menit hingga matang. Angkat dinginkan. Siap digoreng dengan api kecil sebentar. Anakku suka digoreng agak lama dengan kuning telur kocok sisa adonan tadi. Nanti jadi ada kayak jala2 crispynya yang gurih & kering. Angkat & tiriskan jika sudah digoreng matang.--Sajikan dengan dipotong diagonal (potong 2 serong). Siap disantap dengan saus sambal botol yang biasanya saya rebus dengan bawang goreng\", \"cabe rawit iris\", \"gula & sedikit air. Selamat mencoba.--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:21:07', '2025-10-26 07:21:47'),
(10, 3, 'Ayam goreng tulang lunak', 'Resep Ayam goreng tulang lunak yang lezat dan mudah dibuat.', '[\"1 kg ayam (dipotong sesuai selera jangan kecil2 ya)--2 batang serai (memarkan)--4 lembar daun jeruk--7 butir bawang putih (haluskan)--1 sdm ketumbar (haluskan)--3 ruas jari Laos (haluskan)--3 ruas jari kunyit (haluskan)--2 butir kemiri (haluskan)--secukupnya Garam--Secukupnya Air (tuk ukep ayam)--Secukupnya Minyak goreng--\"]', '[\"Haluskan bumbu2nya (BaPut, ketumbar, kemiri, kunyit, Laos, garam) hingga halus, sisihkan--Campur kan bumbu halus tadi dengan ayam yg sudah dicuci bersih dan sudah dipotong didalam panci presto\", \"Uleni sampai tercampur rata\", \"--Tambahkan air hingga ayam tenggelam semua\", \"Masukkan serai dan daun jeruk nya kedalam rendaman ayam\", \"Tutup panci presto rebus/ ukep presto sampai kurleb 45 menit\", \"Dengan api sedang\", \"--Setelah proses ukep presto selesai, tunggu suhu dingin ruang\", \"Lalu goreng ayam dengan minyak goreng api sedang sampai ayam berwarna kecoklatan\", \"--Matang dan sajikan ayam selagi hangat bersama nasi putih, sambal dgn perasan jeruk nipis, lalapan\", \"--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:24:27', '2025-10-26 07:21:51'),
(11, 3, 'Ikan Kembung Bakar Teflon', 'Resep Ikan Kembung Bakar Teflon yang lezat dan mudah dibuat.', '[\"1/2 kg ikan kembung sate, bersihkan--1 buah jeruk sate/jeruk kunci--1 sdm garam halus--2 sdt lada bubuk--1 sdm ketumbar bubuk--\"]', '[\"Kucuri ikan dengan jeruk, diamkan 5 menit--Lumuri ikan dg garam, lada dan ketumbar bubuk, simpan d kulkas sekitar 30 menit--Bakar d atas teflon yg d olesi mentega/ butter tipis dengan api yg kecil--Balik untuk membakar kedua bagian ikan\", \"--Sajikan dengan cocolan sambal kecap--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:25:16', '2025-10-26 07:21:55'),
(12, 3, 'Rabeg Kambing', 'Resep Rabeg Kambing yang lezat dan mudah dibuat.', '[\"1 kg daging kambing bagian paha beserta tulang--Bumbu halus--100 gr cabe--7 siung bawang putih--8 siung bawang mersh--7 butir kemiri--2 cm jahe--1 ruas kunyit--1/2 sdt klabet--1/2 sdt jinten--Bumbu lain--1 batang sereh--3 lembar daun salam--2 cm kayu manis--5 butir cengjeh--3 kapulaga hijau--5 sdm kecap manis--Garam--Gula--Pelengkap--Acar--Bawang goreng--\"]', '[\"Tumis bumbu halus hingga harum. Masukan bumbu bumbu lain kecuali kecap.--Masukan daging kambing. Aduk sampai daging terbalut bumbu dan berubah warna--Pindahkan tumisan daging ke panci. Tambahkan air. Beri kecap gulz dan garam. Masak sampai daging empuk dan air tinggal sedikit--Sajikan dengan nasi dan pelebgkap--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:25:30', '2025-10-26 07:21:59'),
(13, 3, 'Rendang Sapi', 'Resep Rendang Sapi yang lezat dan mudah dibuat.', '[\"1 kg daging sapi--4 butir kelapa untuk 2L santan kental--Bumbu dihaluskan :--200 gram Cabe merah--1/2 cm Jahe--2 sdm ketumbar sangrai--10 butir kapulaga--10 butir cengkeh--2 butir bunga lawang--1 buah biji pala--13 butir bawang merah--7 butir bawang putih--25 buah cabe rawit (jika suka pedas)--3 cm lengkuas dipipihkan--2 buah serai dipipihkan--3 lembar daun jeruk--3 lembar daun salam--2 lembar daun kunyit--1/2 sdm garam (tambahkan jika kurang asin)--1 sdt kaldu bubuk (optional)--\"]', '[\"Potong daging berlawanan dari seratnya--Peras santan dari 4 butir kelapa tua menjadi 2 liter--Haluskan bumbu--Tuang santan dan bumbu ke dalam wajan anti lengket, masak hingga mendidih--Kemudian masukkan daging, garam dan kaldu\", \"Masak dengan api kecil sambil sesekali di aduk agar tidak gosong--Sampai tahap ini dinamakan kalio, teruskan memasak hingga daging mulai mengering dan kehitaman--Rendang siap dinikmati--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:25:45', '2025-10-26 07:22:03'),
(14, 3, 'Batagor (Bakso Tahu Goreng) ala rumahan anti gagal', 'Resep Batagor (Bakso Tahu Goreng) ala rumahan anti gagal yang lezat dan mudah dibuat.', '[\"Bahan batagor--3 kotak ukuran kecil tahu kuning--5 sdm tepung terigu--1/2 sdm mentega (supaya kriuk2)--2 buah bawang putih--Garam (secukupnya)--Merica (secukupnya)--Air putih matang (secukupnya)--Bahan saus--2 sdm saus cabe (saya pake indofo**)--1/2 siung bawang putih--2 cabe rawit (sesuai selera ya karna suami saya gak suka pedas)--Kecap (secukupnya)--Air matang (secukupnya)--\"]', '[\"Haluskan tahu kuning, bawang putih, garam merica (hingga lumat)--Tambahkan mentega, tepung terigu kedalam adonan tahu dan beri air hingga lembek tapi jangan sampai encer ya\", \"Koreksi rasa--Siapkan wajan dengan minyak, goreng adonan yang telah dibuat\", \"(Saya pake 2 sendok buat bikin kecil2)\", \"Angkat dan tata di mangkuk\", \"--Haluskan bawang putih, cabai rawit--Tumis bawang putih dan cabai yang sudah di haluskan hingga harus\", \"--Masukkan saus, kecap, garam\", \"Koreksi rasa--Batagor dan sausnya siap untuk dinikmati 🤤--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:26:04', '2025-10-26 07:22:07'),
(15, 3, 'Terik ayam tempe telor', 'Resep Terik ayam tempe telor yang lezat dan mudah dibuat.', '[\"4 sayap ayam--4 buah telor--1/2 papan tempe--Bumbu halus:--5 siung bawang merah--4 siung bawang putih--2 butir kemiri--1/2 sdm ketumbar--secukupnya Garam--Bumbu lain:--2 cm lengkuas--1 lembar daun salam--1 batang sereh--secukupnya Gula jawa--Bawang goreng--Minyak untuk menumis--1 bungkus santan kara--secukupnya Air--5 buah cabe rawit--\"]', '[\"Cuci sayap ayam kemudian sisihkan. Rebus telur dan kupas. Sisihkan. Potong tempe sesuai selera. Sisihkan.--Panaskan minyak langsung di dalam panci. Tumis sereh\", \"daun salam\", \"dan lengkuas. Setelah beberapa saat masukkan bumbu halus. Tumis sampai berbau harum. Kemudian masukkan air dan gula jawa dan cabe rawit. Jika air sudah mendidih\", \"masukkan sayap ayam dan telor. Masukkan tempe belakangan. Tunggu beberapa saat sampai ayam empuk (karena pakai sayap ayam\", \"jadi tidak begitu lama). Jika ayam sudah empuk masukkan santan kara. Tunggu beberapa saat. Sesuaikan rasa. Matikan api\", \"dan taburkan bawang goreng.--Terik siap dinikmati dengan nasi hangat dan kerupun 😁--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:26:19', '2025-10-26 07:22:11'),
(17, 3, 'Bakso Ayam Udang Keto', 'Resep Bakso Ayam Udang Keto yang lezat dan mudah dibuat.', '[\"400 gr ayam giling--250 gr udang kupas--3 telur ayam utuh kocok--30 ml minyak goreng--1/2 sdt baking powder--secukupnya garam, lada--2 sdm bawang goreng haluskan--secukupnya Bawang merah & putih--1/2 sdt minyak wijen--Karageenan secukupnya (pengenyal)--Es batu 8 keping kotak--\"]', '[\"Campur semua bahan dalam food processor\", \"kecuali daun bawang dan minyak goreng.--Tuang dalam baskom lalu taburi daun bawang dan minyak goreng\", \"uleni rata (gunakan sarung tangan). Lalu simpan 2 jam dlm kulkas.--Panaskan minyak di wajan dgn api sedang. Bulatkan adonan dengan menggunakan dua sendok. Lalu goreng sampai matang. Adonan akan mengembang saat di wajan\", \"tapi kembali kempes stelah diangkat.--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:27:23', '2025-10-26 07:22:22'),
(18, 4, 'Ayam cabai kawin', 'Resep Ayam cabai kawin yang lezat dan mudah dibuat.', '[\"1/4 kg ayam--3 buah cabai hijau besar--7 buah cabai merah rawit--3 siung bawang putih--2 siung bawang merah--secukupnya Gula--secukupnya Garam--1/4 buah tomat merah--secukupnya Air--secukupnya Minyak goreng--\"]', '[\"Panaskan minyak di dalam wajan. Setelah minyak panas masukkan ayam yang sudah dipotong dadu. Goreng hingga matang. Lalu tiriskan.--Haluskan bawang putih\", \"bawang merah\", \"cabai hijau dan merah\", \"tomat.--Panaskan minyak didalam wajan. Setelah minyak panas\", \"masukkan bumbu yang sudah halus. Tunggu sampai wangi. Masukkan ayam yang sudah di goreng. Tambahkan air\", \"gula dan garam. Tunggu sampai bumbu meresap di ayam. Sajikan.--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:33:32', '2025-10-26 07:22:27'),
(19, 4, 'Mujaer asam pedas manis', 'Resep Mujaer asam pedas manis yang lezat dan mudah dibuat.', '[\"1/2 kg ikan mujaer (stok gurame habis)--2 buah Wortel potong korek api--1 siung bawang bombai--7 siung bawang putih--13 cabai rawit--4 siung bawang merah--2 sendok makan Saos tomat--2 sendok makan Saos tiram--Garam--Merica--100 ml air matang--Jahe--Jeruk nipis--\"]', '[\"Bersihkan ikan sampai benar-benar bersih. Kerat-kerat dagingnya. Beri garam\", \"merica\", \"potongan jahe perasan jeruk nipis\", \"dan 2 siung bawang uleg. Diamkan 10 menit--Uleg 5 bawang putih yang tersisa. Potong bawang bombai dan dan bawang merah dan cabai rawit.--Goreng ikan mujaer hingga matang\", \"lalu sisihkan.--Tumis bawang merah dan bawang bombai hingga wangi\", \"masukkan bawang putih uleg. Tumis hingga wangi.--Masukkan wortel\", \"tambahkan air. Aduk-aduk. Tunggu hingga wortel hampir matang.--Masukkan saos tomat dan saos tiram. Beri perasan jeruk nipis. Aduk-aduk Koreksi rasa. Jika sudah pas\", \"masukkan ikan\", \"aduk sebentar saja.--Siap disajikan.--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:34:04', '2025-10-26 07:22:31'),
(20, 4, 'Gulai kambing', 'Resep Gulai kambing yang lezat dan mudah dibuat.', '[\"500 gram daging kambing--Bumbu :--sesuai selera Santan--2 btng serai,memarkan--5 lembar Daun Salam,jeruk--sesuai selera Cabe--2 btr Cengkeh--Bumbu halus:--5 siung bawng merah--4 siung bwng putih--5 butir kemiri--2 cm kunyit--2 cm jahe--2 cm lengkuas--1 sdt ketumbar--Mrica,pala,garam,gula seckpnya--\"]', '[\"Rebus dgng kambing dgn jahe krng lbh 20 menit lalu buang airnya,cc bersih--Tumis bumbu yg sudah dihaluskn,,masukn semua bahan Dan bumbu--Rebus sampe empuk--Terakhir masukn santan,tes rasa,,,--Sajikan dimangkok,yumm ✌️--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:34:15', '2025-10-26 07:22:37'),
(21, 4, 'Tongseng daging sapi pedas', 'Resep Tongseng daging sapi pedas yang lezat dan mudah dibuat.', '[\"300 gr daging sapi (iris tipis panjang)--Bumbu halus :--7 siung bawang merah--4 siung bawang putih--2 cm jahe--1,5 cm kunyit--3 bh kemiri--5 cabe merah keriting--Bumbu pelengkap :--2 batang sereh--2 daun salam--2 daun jeruk--5 bh cabe rawit merah utuh--Toping--1 bh tomat besar (iris sedang)--Kol secukupnya menurut selera (iris sedang)--\"]', '[\"Tumis bumbu halus dengan sedikit minyak\", \"tumis sampai matang dan mengental. Lalu masukkan daging sapi yang telah dipotong2 tadi. Lalu masukkan bumbu pelengkap. Tumis daging dengan bumbu sampai daging sapi berubah warna dan agak mengeluarkan minyak.--Masukkan air secukupnya sekitar 1\", \"5 liter biarkan mendidih sebentar lalu kecilkan api dan tutup panci. Tunggu sampai daging empuk..--Setelah daging empuk masukkan kol dan tomat\", \"masak sebentar saja agar tekstur kol agak crunchy.. angkat lalu sajikan--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'pending', '', '2025-10-23 15:34:27', '2025-10-23 15:34:27'),
(22, 4, 'Sayur Sutace (sop tahu ceker)', 'Resep Sayur Sutace (sop tahu ceker) yang lezat dan mudah dibuat.', '[\"1/2 kg ceker ayam--4 buah wortel ukuran sedang--6 buah wortel--sesuai selera Kol--sesuai selera Daun bawang--sesuai selera Seledri--4 buah tahu kuning--Garam--Gula putih--Penyedap rasa (ma**ko)--Lada bubuk--Jeruk nipis--2 siung Bawang Putih--2 siung Bawang merah--\"]', '[\"Bersihkan ceker, lau rebut selama 3-4 menit pakai jeruk nipis biar gk terlalu bau--Tiriskan ceker--Potong2 tahu, lalu goreng set mateng, tiriskan--Potong2 sayurr, cuci hingga bersihh\", \"--Haluskan bawang merah dan bawang putih\", \"--Panaskan sedikit minyak masukan bawang yg sudah di haluskan, oseng2 hingga tercium wangi--Masukan sayurr, tampahkan air, garam, gula putih dan penyedap rasa secukupnya\", \"--Tunggu hingga sayur hampir matang, masukan ceker, dan tahu\", \"--Tutup wajan hingga kaldu meresap kedalam ceker dan tahu\", \"--Setelah matang, tuang ke mangkok, siap di santappp 😊😊--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'pending', '', '2025-10-23 15:34:45', '2025-10-23 15:34:45'),
(23, 4, 'Telur Kornet', 'Resep Telur Kornet yang lezat dan mudah dibuat.', '[\"1/2 kaleng kornet--2-3 buah bawang prei--2 buah telur--1 sdm tepung terigu--secukupnya Cabe rawit--2 jumput Garam--1/2 sdt Merica--1/2 sdm Saus tiram--\"]', '[\"Campur kornet\", \"telur\", \"bawang prei dan cabe yg sudah di iris tipis. Tambahkan garam merica saos tiram secukupnya dan tepung terigu. Goreng 1 sdm dulu setelah matang baru cek rasa. Dirasa sudah pas. Adonan siap di goreng.😊--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'pending', '', '2025-10-23 15:34:58', '2025-10-23 15:34:58'),
(24, 4, 'Penyet Tempe Sambel Korek Kemangi', 'Resep Penyet Tempe Sambel Korek Kemangi yang lezat dan mudah dibuat.', '[\"2 buah tempe--1 genggam Daun kemangi--15 cabe rawit--2 siung bawang putih--Gula--Garam--\"]', '[\"Iris tempe jd beberapa potong, bumbui dg bawang putih dan garam\", \"Goreng sampai kering--Ulek cabe rawit dan 1 siung bawang putih ukuran besar tambah gula garam sesuai selera--Setelah halus tuang minyak goreng panas bekas goreng tempe td--Penyet tempe diatas sambel dan beri daun kemangi--Tempe penyet siap disantap 😋😋😋--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'pending', '', '2025-10-23 15:35:11', '2025-10-23 15:35:11'),
(25, 4, 'Udang ala pop corn', 'Resep Udang ala pop corn yang lezat dan mudah dibuat.', '[\"1/4 kg udang basah ukuran sedang--1 bungkus kobe tepung ayam super crispy--secukupnya Air matang--Minyak untuk menggoreng--\"]', '[\"Buang kepala dan cangkang udang.--Cuci bersih udang.--Tepung bumbu dibagi jadi 2 adonan. Adonan basah dan adonan kering.--Masukkan udang ke dalam adonan tepung bumbu kering\", \"gulirkan ke dlm tepung sambil ditekan2. Lalu masukkan ke dalam Tepung bumbu adonan basah gulirkan lagi sambil di tekan2\", \"lalu masukkan lagi di adonan tepung bumbu kering Sambil ditekan2 agar tepung menempel sempurna.--Panaskan minyak goreng. Lalu goreng udang sampai kuning keemasan.--Angkat dan tiriskan.--Sajikan selagi hangat 😊--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'rejected', 'Instruksi kurang jelas.', '2025-10-23 15:35:58', '2025-10-26 06:04:56'),
(26, 4, 'Ayam Geprek', 'Resep Ayam Geprek yang lezat dan mudah dibuat.', '[\"250 gr daging ayam (saya pakai fillet)\", \"Secukupnya gula dan garam\", \"50-100 gr tepung ayam serbaguna\", \"Secukupnya lalapan (kemangi,kol,timun)\", \"Secukupnya minyak panas\", \"❤sambal korek\", \"Secukupnya cabe rawit merah dan bwg putih\"]', '[\"Goreng ayam seperti ayam krispi\", \"Ulek semua bahan sambal kemudian campur dengan minyak panas bekas goreng ayam\", \"Geprek ayam kemudian campur dengan sambal,sajikan dengan lalapan ❤\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 16:23:16', '2025-10-26 06:02:33'),
(29, 2, 'semur daging', 'semur daging ', '[\"1/2 daging sapi\", \"royco 1kg\"]', '[\"semur daging\", \"sampai mateng\"]', 'http://localhost:8080/assets/semur-daging-1.jpg', 'pending', '', '2025-11-09 06:02:17', '2025-11-09 06:02:17');

-- --------------------------------------------------------

--
-- Table structure for table `menu_tags`
--

CREATE TABLE `menu_tags` (
  `menu_id` bigint UNSIGNED NOT NULL,
  `tag_id` bigint UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `menu_tags`
--

INSERT INTO `menu_tags` (`menu_id`, `tag_id`) VALUES
(26, 1),
(29, 4),
(26, 20),
(26, 50),
(29, 50),
(29, 53);

-- --------------------------------------------------------

--
-- Table structure for table `menu_votes`
--

CREATE TABLE `menu_votes` (
  `vote_id` bigint UNSIGNED NOT NULL,
  `menu_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `likes_count` int NOT NULL DEFAULT '0',
  `dislikes_count` int NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `menu_votes`
--

INSERT INTO `menu_votes` (`vote_id`, `menu_id`, `user_id`, `likes_count`, `dislikes_count`, `created_at`, `updated_at`) VALUES
(5, 3, 6, 1, 0, '2025-10-26 07:09:15', '2025-11-09 07:17:49'),
(6, 4, 6, 1, 0, '2025-10-26 07:09:21', '2025-11-09 07:17:49'),
(7, 1, 3, 1, 0, '2025-10-26 07:10:22', '2025-11-09 07:17:49'),
(8, 2, 3, 1, 0, '2025-10-26 07:10:26', '2025-11-09 07:17:49'),
(9, 3, 3, 1, 0, '2025-10-26 07:10:30', '2025-11-09 07:17:49'),
(10, 3, 4, 1, 0, '2025-10-26 07:11:18', '2025-11-09 07:17:49'),
(11, 2, 4, 1, 0, '2025-10-26 07:11:22', '2025-11-09 07:17:49'),
(12, 1, 4, 1, 0, '2025-10-26 07:11:26', '2025-11-09 07:17:49'),
(13, 1, 7, 1, 0, '2025-10-26 14:14:10', '2025-11-09 07:17:49'),
(14, 5, 2, 0, 1, '2025-11-08 09:54:39', '2025-11-09 07:17:49'),
(20, 11, 2, 1, 0, '2025-11-08 10:10:58', '2025-11-09 07:17:49'),
(21, 10, 2, 1, 0, '2025-11-08 10:11:08', '2025-11-09 07:17:49');

-- --------------------------------------------------------

--
-- Table structure for table `password_resets`
--

CREATE TABLE `password_resets` (
  `reset_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `roles`
--

CREATE TABLE `roles` (
  `role_id` bigint UNSIGNED NOT NULL,
  `role_name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `roles`
--

INSERT INTO `roles` (`role_id`, `role_name`) VALUES
(1, 'admin'),
(2, 'editor'),
(3, 'member');

-- --------------------------------------------------------

--
-- Table structure for table `tags`
--

CREATE TABLE `tags` (
  `tag_id` bigint UNSIGNED NOT NULL,
  `category_id` bigint UNSIGNED NOT NULL,
  `tag_name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `tags`
--

INSERT INTO `tags` (`tag_id`, `category_id`, `tag_name`) VALUES
(1, 1, 'Ayam'),
(2, 1, 'Ikan'),
(3, 1, 'Kambing'),
(4, 1, 'Sapi'),
(5, 1, 'Tahu'),
(6, 1, 'Telur'),
(7, 1, 'Tempe'),
(8, 1, 'Udang'),
(9, 1, 'Nasi'),
(10, 1, 'Mie'),
(11, 1, 'Sayuran'),
(20, 2, 'Goreng'),
(21, 2, 'Bakar'),
(22, 2, 'Rebus'),
(23, 2, 'Tumis'),
(24, 2, 'Kuah'),
(25, 2, 'Panggang'),
(26, 2, 'Sate'),
(27, 2, 'Soto'),
(28, 2, 'Sup'),
(29, 2, 'Gulai'),
(40, 3, 'Sarapan'),
(41, 3, 'Makan Siang'),
(42, 3, 'Makan Malam'),
(43, 3, 'Dessert'),
(44, 3, 'Camilan'),
(50, 4, 'Pedas'),
(51, 4, 'Manis'),
(52, 4, 'Asam'),
(53, 4, 'Gurih'),
(54, 4, 'Asin');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `user_id` bigint UNSIGNED NOT NULL,
  `role_id` bigint UNSIGNED NOT NULL,
  `status_user` enum('active','inactive') NOT NULL DEFAULT 'active',
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`user_id`, `role_id`, `status_user`, `name`, `email`, `password`, `created_at`, `updated_at`) VALUES
(1, 1, 'active', 'Admin Utama', 'admin@gmail.com', '$2a$10$ghzwfho7MvLgg247iumqve7m.XZYZHSYeiVIigtB3ZHRufvF/9Elm', '2025-10-23 13:55:32', '2025-10-23 13:55:32'),
(2, 3, 'active', 'Budi Santoso', 'budi@gmail.com', '$2a$10$4tSbPmmnpscN2Gfyse9GieU8aTPeQMrQC/EqxCSTG/urVFANwsMCq', '2025-10-23 13:56:42', '2025-11-08 10:11:07'),
(3, 3, 'active', 'Ani Lestari', 'ani@gmail.com', '$2a$10$c5pLKPfMoO2h/5mYVq/wieFfeseTHThLFz7W9F0j3EQi4SqOrmMUW', '2025-10-23 13:57:06', '2025-10-26 07:10:07'),
(4, 2, 'active', 'Candra Wijaya', 'candra@gmail.com', '$2a$10$LcnkpHDgymYKSPI2Ri5qoetb2l4bRMrppSArtF/q6GORwTY7TmzYO', '2025-10-23 13:57:25', '2025-10-26 07:11:10'),
(5, 2, 'active', 'Editor', 'editor@gmail.com', '$2a$10$FX8E3Pplm58noZ44FLGH.ePT5ecVcjZmypNRrY3rxtbVmVQSeFGc.', '2025-10-23 13:58:08', '2025-10-23 14:30:02'),
(6, 3, 'active', 'Asep Wiyanto', 'asep@gmail.com', '$2a$10$p1/AQxEl50ml2sAgl4YdperlBm98JXH1pWijGt.r8P6cblDWKlXJe', '2025-10-26 06:25:26', '2025-10-26 07:07:54'),
(7, 3, 'active', 'Udin Wibawa', 'udin@gmail.com', '$2a$10$pwDQRpjTwBYKikhE.z55W.tem6C5ogMHHTLcqPiwqe/DfGam3aRbW', '2025-10-26 14:11:46', '2025-10-26 14:14:30');

-- --------------------------------------------------------

--
-- Table structure for table `user_bookmarks`
--

CREATE TABLE `user_bookmarks` (
  `user_id` bigint UNSIGNED NOT NULL,
  `menu_id` bigint UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `user_bookmarks`
--

INSERT INTO `user_bookmarks` (`user_id`, `menu_id`) VALUES
(3, 1),
(4, 1),
(6, 1),
(7, 1),
(3, 2),
(4, 2),
(6, 2),
(2, 3),
(4, 3),
(6, 3),
(2, 5),
(2, 7),
(2, 10),
(2, 11);

-- --------------------------------------------------------

--
-- Table structure for table `user_profiles`
--

CREATE TABLE `user_profiles` (
  `profile_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `bio` text,
  `profile_picture_url` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `user_profiles`
--

INSERT INTO `user_profiles` (`profile_id`, `user_id`, `bio`, `profile_picture_url`) VALUES
(1, 1, '', ''),
(2, 2, '', ''),
(3, 3, '', ''),
(4, 4, '', ''),
(5, 5, '', ''),
(6, 6, '', ''),
(7, 7, '', '');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`category_id`);

--
-- Indexes for table `comments`
--
ALTER TABLE `comments`
  ADD PRIMARY KEY (`comment_id`),
  ADD KEY `fk_comments_menu` (`menu_id`),
  ADD KEY `fk_comments_user` (`user_id`);

--
-- Indexes for table `log_activity`
--
ALTER TABLE `log_activity`
  ADD PRIMARY KEY (`activity_id`),
  ADD KEY `fk_log_user` (`user_id`);

--
-- Indexes for table `menus`
--
ALTER TABLE `menus`
  ADD PRIMARY KEY (`menu_id`),
  ADD KEY `fk_menus_user` (`user_id`);

--
-- Indexes for table `menu_tags`
--
ALTER TABLE `menu_tags`
  ADD PRIMARY KEY (`menu_id`,`tag_id`),
  ADD KEY `fk_menutags_tag` (`tag_id`);

--
-- Indexes for table `menu_votes`
--
ALTER TABLE `menu_votes`
  ADD PRIMARY KEY (`vote_id`),
  ADD UNIQUE KEY `user_menu_vote_unique` (`user_id`,`menu_id`),
  ADD KEY `fk_votes_menu` (`menu_id`);

--
-- Indexes for table `password_resets`
--
ALTER TABLE `password_resets`
  ADD PRIMARY KEY (`reset_id`),
  ADD UNIQUE KEY `token_unique` (`token`),
  ADD KEY `fk_resets_user` (`user_id`);

--
-- Indexes for table `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`role_id`),
  ADD UNIQUE KEY `role_name_unique` (`role_name`);

--
-- Indexes for table `tags`
--
ALTER TABLE `tags`
  ADD PRIMARY KEY (`tag_id`),
  ADD UNIQUE KEY `tag_name_unique` (`tag_name`),
  ADD KEY `fk_tags_category` (`category_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `email_unique` (`email`),
  ADD KEY `fk_users_role` (`role_id`);

--
-- Indexes for table `user_bookmarks`
--
ALTER TABLE `user_bookmarks`
  ADD PRIMARY KEY (`user_id`,`menu_id`),
  ADD KEY `fk_bookmarks_menu` (`menu_id`);

--
-- Indexes for table `user_profiles`
--
ALTER TABLE `user_profiles`
  ADD PRIMARY KEY (`profile_id`),
  ADD UNIQUE KEY `user_id_unique` (`user_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `categories`
--
ALTER TABLE `categories`
  MODIFY `category_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `comments`
--
ALTER TABLE `comments`
  MODIFY `comment_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `log_activity`
--
ALTER TABLE `log_activity`
  MODIFY `activity_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=171;

--
-- AUTO_INCREMENT for table `menus`
--
ALTER TABLE `menus`
  MODIFY `menu_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=30;

--
-- AUTO_INCREMENT for table `menu_votes`
--
ALTER TABLE `menu_votes`
  MODIFY `vote_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=22;

--
-- AUTO_INCREMENT for table `password_resets`
--
ALTER TABLE `password_resets`
  MODIFY `reset_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `roles`
--
ALTER TABLE `roles`
  MODIFY `role_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `tags`
--
ALTER TABLE `tags`
  MODIFY `tag_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=56;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `user_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `user_profiles`
--
ALTER TABLE `user_profiles`
  MODIFY `profile_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `comments`
--
ALTER TABLE `comments`
  ADD CONSTRAINT `fk_comments_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

--
-- Constraints for table `log_activity`
--
ALTER TABLE `log_activity`
  ADD CONSTRAINT `fk_log_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

--
-- Constraints for table `menus`
--
ALTER TABLE `menus`
  ADD CONSTRAINT `fk_menus_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

--
-- Constraints for table `menu_tags`
--
ALTER TABLE `menu_tags`
  ADD CONSTRAINT `fk_menutags_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_menutags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`tag_id`) ON DELETE CASCADE;

--
-- Constraints for table `menu_votes`
--
ALTER TABLE `menu_votes`
  ADD CONSTRAINT `fk_votes_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_votes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

--
-- Constraints for table `password_resets`
--
ALTER TABLE `password_resets`
  ADD CONSTRAINT `fk_resets_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

--
-- Constraints for table `tags`
--
ALTER TABLE `tags`
  ADD CONSTRAINT `fk_tags_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`category_id`);

--
-- Constraints for table `users`
--
ALTER TABLE `users`
  ADD CONSTRAINT `fk_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`role_id`);

--
-- Constraints for table `user_bookmarks`
--
ALTER TABLE `user_bookmarks`
  ADD CONSTRAINT `fk_bookmarks_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_bookmarks_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

--
-- Constraints for table `user_profiles`
--
ALTER TABLE `user_profiles`
  ADD CONSTRAINT `fk_profiles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
