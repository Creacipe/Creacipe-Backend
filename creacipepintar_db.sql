-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Nov 22, 2025 at 01:51 PM
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
  `parent_id` bigint UNSIGNED DEFAULT NULL,
  `comment_text` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `comments`
--

INSERT INTO `comments` (`comment_id`, `menu_id`, `user_id`, `parent_id`, `comment_text`, `created_at`) VALUES
(1, 1, 13, NULL, 'mantap mudah di pahami resep nya', '2025-11-22 07:14:10'),
(2, 1, 2, 1, 'terima kasih banyak semoga bermanfaat', '2025-11-22 07:53:20');

-- --------------------------------------------------------

--
-- Table structure for table `email_verifications`
--

CREATE TABLE `email_verifications` (
  `verification_id` int NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `new_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `verification_code` varchar(6) COLLATE utf8mb4_unicode_ci NOT NULL,
  `expires_at` datetime NOT NULL,
  `is_used` tinyint(1) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `email_verifications`
--

INSERT INTO `email_verifications` (`verification_id`, `user_id`, `new_email`, `verification_code`, `expires_at`, `is_used`, `created_at`) VALUES
(5, 1, 'maiys031104@gmail.com', '327346', '2025-11-17 09:34:59', 1, '2025-11-17 09:24:59'),
(7, 1, 'damnihbos4@gmail.com', '130797', '2025-11-17 10:10:29', 1, '2025-11-17 10:00:29'),
(8, 1, 'maiys031104@gmail.com', '413089', '2025-11-21 17:34:59', 0, '2025-11-21 17:24:59');

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
(170, 2, 'USER_LOGIN', 2, 'users', '2025-11-09 06:49:39'),
(171, 2, 'USER_LOGIN', 2, 'users', '2025-11-09 07:38:20'),
(172, 2, 'UNLIKE_MENU', 11, 'menus', '2025-11-09 07:48:34'),
(173, 2, 'LIKE_MENU', 11, 'menus', '2025-11-09 07:48:35'),
(174, 2, 'LIKE_MENU', 1, 'menus', '2025-11-09 07:48:41'),
(175, 2, 'LIKE_MENU', 7, 'menus', '2025-11-09 07:59:54'),
(176, 2, 'BOOKMARK_MENU', 8, 'menus', '2025-11-09 08:28:42'),
(177, 2, 'USER_LOGIN', 2, 'users', '2025-11-09 09:08:24'),
(178, 2, 'UNBOOKMARK_MENU', 7, 'menus', '2025-11-09 10:28:54'),
(179, 2, 'BOOKMARK_MENU', 7, 'menus', '2025-11-09 10:28:55'),
(180, 2, 'UNBOOKMARK_MENU', 7, 'menus', '2025-11-09 10:28:56'),
(181, 2, 'BOOKMARK_MENU', 7, 'menus', '2025-11-09 10:28:56'),
(182, 2, 'UNLIKE_MENU', 7, 'menus', '2025-11-09 10:28:57'),
(183, 2, 'LIKE_MENU', 7, 'menus', '2025-11-09 10:28:58'),
(184, 2, 'UNLIKE_MENU', 7, 'menus', '2025-11-09 10:28:58'),
(185, 2, 'LIKE_MENU', 7, 'menus', '2025-11-09 10:28:59'),
(186, 2, 'LIKE_MENU', 6, 'menus', '2025-11-09 11:08:57'),
(187, 2, 'LIKE_MENU', 15, 'menus', '2025-11-09 11:15:57'),
(188, 2, 'UNLIKE_MENU', 15, 'menus', '2025-11-09 11:15:58'),
(189, 2, 'LIKE_MENU', 15, 'menus', '2025-11-09 11:15:59'),
(190, 2, 'BOOKMARK_MENU', 15, 'menus', '2025-11-09 11:16:00'),
(191, 2, 'UNBOOKMARK_MENU', 15, 'menus', '2025-11-09 11:16:01'),
(192, 2, 'BOOKMARK_MENU', 15, 'menus', '2025-11-09 11:16:02'),
(193, 2, 'UNLIKE_MENU', 15, 'menus', '2025-11-09 11:58:02'),
(194, 2, 'UNBOOKMARK_MENU', 15, 'menus', '2025-11-09 11:58:03'),
(195, 2, 'BOOKMARK_MENU', 15, 'menus', '2025-11-09 11:58:04'),
(196, 2, 'LIKE_MENU', 15, 'menus', '2025-11-09 11:58:05'),
(197, 2, 'UNDISLIKE_MENU', 5, 'menus', '2025-11-09 11:59:11'),
(198, 2, 'UNBOOKMARK_MENU', 5, 'menus', '2025-11-09 11:59:13'),
(199, 2, 'UNLIKE_MENU', 15, 'menus', '2025-11-09 12:05:34'),
(200, 2, 'LIKE_MENU', 15, 'menus', '2025-11-09 12:05:35'),
(201, 2, 'UNLIKE_MENU', 15, 'menus', '2025-11-09 12:05:35'),
(202, 2, 'LIKE_MENU', 15, 'menus', '2025-11-09 12:05:36'),
(203, 2, 'LIKE_MENU', 4, 'menus', '2025-11-09 12:07:25'),
(204, 2, 'DISLIKE_MENU', 12, 'menus', '2025-11-09 12:07:29'),
(205, 2, 'LIKE_MENU', 12, 'menus', '2025-11-09 12:11:39'),
(206, 2, 'UNLIKE_MENU', 12, 'menus', '2025-11-09 12:11:40'),
(207, 2, 'LIKE_MENU', 12, 'menus', '2025-11-09 12:11:41'),
(208, 2, 'UNBOOKMARK_MENU', 10, 'menus', '2025-11-09 12:36:26'),
(209, 2, 'BOOKMARK_MENU', 10, 'menus', '2025-11-09 12:45:45'),
(210, 2, 'UNBOOKMARK_MENU', 10, 'menus', '2025-11-09 12:45:54'),
(211, 2, 'UPDATE_MENU', 29, 'menus', '2025-11-09 14:30:48'),
(212, 2, 'CREATE_MENU', 30, 'menus', '2025-11-09 14:32:21'),
(213, 2, 'DELETE_MENU', 29, 'menus', '2025-11-09 14:34:31'),
(214, 2, 'UPDATE_MENU', 30, 'menus', '2025-11-09 14:53:34'),
(215, 2, 'CREATE_MENU', 31, 'menus', '2025-11-09 14:54:22'),
(216, 2, 'CREATE_MENU', 32, 'menus', '2025-11-09 15:21:13'),
(217, 2, 'UPDATE_MENU', 32, 'menus', '2025-11-09 15:23:09'),
(218, 2, 'UPDATE_MENU', 31, 'menus', '2025-11-09 15:56:12'),
(219, 2, 'UPDATE_MENU', 32, 'menus', '2025-11-09 16:01:21'),
(220, 2, 'DELETE_MENU', 30, 'menus', '2025-11-09 16:01:46'),
(221, 2, 'USER_LOGIN', 2, 'users', '2025-11-09 16:08:33'),
(222, 2, 'UPDATE_MENU', 9, 'menus', '2025-11-09 16:09:50'),
(223, 2, 'UPDATE_MENU', 2, 'menus', '2025-11-09 16:10:13'),
(224, 2, 'UPDATE_MENU', 1, 'menus', '2025-11-09 16:10:24'),
(225, 2, 'UPDATE_MENU', 3, 'menus', '2025-11-09 16:10:35'),
(226, 2, 'UPDATE_MENU', 4, 'menus', '2025-11-09 16:10:44'),
(227, 2, 'UPDATE_MENU', 5, 'menus', '2025-11-09 16:10:55'),
(228, 2, 'UPDATE_MENU', 6, 'menus', '2025-11-09 16:11:14'),
(229, 2, 'UPDATE_MENU', 7, 'menus', '2025-11-09 16:11:22'),
(230, 2, 'UPDATE_MENU', 8, 'menus', '2025-11-09 16:11:32'),
(231, 2, 'BOOKMARK_MENU', 2, 'menus', '2025-11-14 06:29:54'),
(232, 2, 'UNBOOKMARK_MENU', 2, 'menus', '2025-11-14 06:29:55'),
(233, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:09:09'),
(234, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:20:50'),
(235, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:22:32'),
(236, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:24:37'),
(237, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:28:59'),
(238, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:29:59'),
(239, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:36:17'),
(240, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:42:14'),
(241, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:42:56'),
(242, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:46:59'),
(243, 1, 'USER_LOGIN', 1, 'users', '2025-11-14 12:51:08'),
(244, 1, 'USER_LOGIN', 1, 'users', '2025-11-15 06:48:07'),
(245, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-15 06:50:21'),
(246, 1, 'USER_LOGIN', 1, 'users', '2025-11-15 06:54:38'),
(247, 5, 'USER_LOGIN', 5, 'users', '2025-11-15 06:56:11'),
(248, 2, 'USER_LOGIN', 2, 'users', '2025-11-15 12:41:24'),
(249, 1, 'USER_LOGIN', 1, 'users', '2025-11-15 12:46:24'),
(250, 1, 'DELETE_MENU', 32, 'menus', '2025-11-15 14:34:23'),
(251, 5, 'USER_LOGIN', 5, 'users', '2025-11-15 14:36:51'),
(252, 5, 'APPROVE_MENU', 23, 'menus', '2025-11-15 14:48:15'),
(253, 5, 'REJECT_MENU', 31, 'menus', '2025-11-15 15:01:33'),
(254, 5, 'UPDATE_TAG', 2, 'tags', '2025-11-15 15:45:33'),
(255, 5, 'UPDATE_TAG', 2, 'tags', '2025-11-15 15:45:39'),
(256, 5, 'CREATE_TAG', 56, 'tags', '2025-11-15 15:45:50'),
(257, 5, 'DELETE_TAG', 56, 'tags', '2025-11-15 15:45:58'),
(258, 5, 'UPDATE_CATEGORY', 1, 'categories', '2025-11-15 15:52:40'),
(259, 5, 'CREATE_CATEGORY', 6, 'categories', '2025-11-15 15:54:18'),
(260, 5, 'UPDATE_CATEGORY', 6, 'categories', '2025-11-15 15:54:29'),
(261, 5, 'UPDATE_CATEGORY', 6, 'categories', '2025-11-15 15:54:34'),
(262, 5, 'DELETE_CATEGORY', 6, 'categories', '2025-11-15 15:54:39'),
(263, 5, 'CREATE_CATEGORY', 7, 'categories', '2025-11-15 16:11:54'),
(264, 5, 'UPDATE_CATEGORY', 7, 'categories', '2025-11-15 16:11:59'),
(265, 5, 'DELETE_CATEGORY', 7, 'categories', '2025-11-15 16:12:03'),
(266, 1, 'USER_LOGIN', 1, 'users', '2025-11-15 16:12:36'),
(267, 1, 'UPDATE_USER', 5, 'users', '2025-11-15 16:13:09'),
(268, 1, 'UPDATE_USER', 5, 'users', '2025-11-16 08:41:21'),
(269, 1, 'UPDATE_USER', 5, 'users', '2025-11-16 08:43:52'),
(270, 1, 'DELETE_USER', 7, 'users', '2025-11-16 08:47:29'),
(271, 1, 'ADMIN_DEACTIVATE_USER', 6, 'users', '2025-11-16 09:11:08'),
(272, 1, 'ADMIN_ACTIVATE_USER', 6, 'users', '2025-11-16 09:11:18'),
(273, 1, 'CREATE_USER', 8, 'users', '2025-11-16 09:14:34'),
(274, 1, 'UPDATE_USER', 8, 'users', '2025-11-16 09:14:45'),
(275, 1, 'UPDATE_USER_ROLE', 8, 'users', '2025-11-16 09:14:56'),
(276, 1, 'UPDATE_USER_ROLE', 8, 'users', '2025-11-16 09:15:07'),
(277, 1, 'ADMIN_DEACTIVATE_USER', 8, 'users', '2025-11-16 09:15:12'),
(278, 1, 'ADMIN_ACTIVATE_USER', 8, 'users', '2025-11-16 09:20:05'),
(279, 1, 'UPDATE_USER', 8, 'users', '2025-11-16 09:20:16'),
(280, 1, 'CREATE_USER', 9, 'users', '2025-11-16 09:22:04'),
(281, 1, 'CREATE_USER', 10, 'users', '2025-11-16 09:22:33'),
(282, 1, 'CREATE_USER', 11, 'users', '2025-11-16 09:23:03'),
(283, 1, 'CREATE_USER', 12, 'users', '2025-11-16 09:27:06'),
(284, 1, 'ADMIN_DEACTIVATE_USER', 11, 'users', '2025-11-16 09:27:18'),
(285, 1, 'ADMIN_DEACTIVATE_USER', 10, 'users', '2025-11-16 09:27:22'),
(286, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 12:41:17'),
(287, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 12:44:12'),
(288, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 12:44:20'),
(289, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 12:44:29'),
(290, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 12:46:42'),
(291, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 12:53:07'),
(292, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 12:53:18'),
(293, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 13:07:22'),
(294, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 13:07:56'),
(295, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 13:25:16'),
(296, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 13:27:20'),
(297, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 13:32:32'),
(298, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 13:33:28'),
(299, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 13:39:42'),
(300, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 15:51:31'),
(301, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-16 16:00:40'),
(302, 1, 'REQUEST_PASSWORD_RESET', 1, 'users', '2025-11-16 16:08:30'),
(303, 1, 'REQUEST_EMAIL_CHANGE', 1, 'users', '2025-11-16 16:13:23'),
(304, 1, 'REQUEST_PASSWORD_RESET', 1, 'users', '2025-11-16 16:19:20'),
(305, 1, 'REQUEST_EMAIL_CHANGE', 1, 'users', '2025-11-16 16:19:55'),
(306, 1, 'REQUEST_EMAIL_CHANGE', 1, 'users', '2025-11-16 16:25:14'),
(307, 1, 'REQUEST_PASSWORD_RESET', 1, 'users', '2025-11-16 16:27:01'),
(308, 1, 'REQUEST_PASSWORD_RESET', 1, 'users', '2025-11-16 16:37:17'),
(309, 1, 'REQUEST_EMAIL_CHANGE', 1, 'users', '2025-11-17 02:25:00'),
(310, 1, 'EMAIL_CHANGE_SUCCESS', 1, 'users', '2025-11-17 02:26:39'),
(311, 1, 'REQUEST_EMAIL_CHANGE', 1, 'users', '2025-11-17 03:00:30'),
(312, 1, 'EMAIL_CHANGE_SUCCESS', 1, 'users', '2025-11-17 03:01:21'),
(313, 1, 'REQUEST_PASSWORD_RESET', 1, 'users', '2025-11-17 03:01:42'),
(314, 1, 'USER_LOGIN', 1, 'users', '2025-11-17 03:17:58'),
(315, 1, 'UPDATE_PROFILE', 1, 'users', '2025-11-17 03:18:25'),
(316, 1, 'REQUEST_PASSWORD_RESET', 1, 'users', '2025-11-17 03:19:01'),
(317, 1, 'PASSWORD_RESET_SUCCESS', 1, 'users', '2025-11-17 03:19:15'),
(318, 1, 'REQUEST_PASSWORD_RESET', 1, 'users', '2025-11-21 09:30:31'),
(319, 1, 'USER_LOGIN', 1, 'users', '2025-11-21 10:22:41'),
(320, 1, 'REQUEST_EMAIL_CHANGE', 1, 'users', '2025-11-21 10:25:00'),
(321, 13, 'USER_REGISTER', 13, 'users', '2025-11-21 10:27:31'),
(322, 13, 'USER_LOGIN', 13, 'users', '2025-11-21 10:27:42'),
(323, 13, 'UPDATE_PROFILE', 13, 'users', '2025-11-21 10:31:37'),
(324, 13, 'FORGOT_PASSWORD_REQUEST', 13, 'users', '2025-11-21 10:32:16'),
(325, 13, 'FORGOT_PASSWORD_SUCCESS', 13, 'users', '2025-11-21 10:36:01'),
(326, 13, 'USER_LOGIN', 13, 'users', '2025-11-21 10:36:12'),
(327, 13, 'LIKE_MENU', 1, 'menus', '2025-11-21 10:38:40'),
(328, 13, 'LIKE_MENU', 15, 'menus', '2025-11-21 10:39:48'),
(329, 2, 'USER_LOGIN', 2, 'users', '2025-11-21 11:35:24'),
(330, 4, 'USER_LOGIN', 4, 'users', '2025-11-21 14:05:27'),
(331, 1, 'USER_LOGIN', 1, 'users', '2025-11-21 14:15:30'),
(332, 1, 'UPDATE_USER_ROLE', 4, 'users', '2025-11-21 14:15:45'),
(333, 2, 'USER_LOGIN', 2, 'users', '2025-11-21 14:16:12'),
(334, 2, 'CREATE_MENU', 33, 'menus', '2025-11-21 14:35:29'),
(335, 5, 'USER_LOGIN', 5, 'users', '2025-11-21 14:37:01'),
(336, 5, 'REJECT_MENU', 33, 'menus', '2025-11-21 14:38:05'),
(337, 4, 'USER_LOGIN', 4, 'users', '2025-11-21 14:38:33'),
(338, 2, 'USER_LOGIN', 2, 'users', '2025-11-21 14:39:28'),
(339, 2, 'BOOKMARK_MENU', 12, 'menus', '2025-11-21 15:50:10'),
(340, 2, 'UNBOOKMARK_MENU', 12, 'menus', '2025-11-21 16:01:11'),
(341, 2, 'BOOKMARK_MENU', 12, 'menus', '2025-11-21 16:16:59'),
(342, 2, 'BOOKMARK_MENU', 6, 'menus', '2025-11-21 16:19:06'),
(343, 2, 'UNBOOKMARK_MENU', 6, 'menus', '2025-11-21 16:19:25'),
(344, 2, 'UNLIKE_MENU', 4, 'menus', '2025-11-22 06:26:34'),
(345, 2, 'BOOKMARK_MENU', 15, 'menus', '2025-11-22 06:27:21'),
(346, 2, 'UNBOOKMARK_MENU', 15, 'menus', '2025-11-22 06:27:23'),
(347, 2, 'UNLIKE_MENU', 15, 'menus', '2025-11-22 06:29:46'),
(348, 2, 'LIKE_MENU', 15, 'menus', '2025-11-22 06:29:48'),
(349, 2, 'UNLIKE_MENU', 15, 'menus', '2025-11-22 06:29:49'),
(350, 2, 'UNLIKE_MENU', 1, 'menus', '2025-11-22 06:44:18'),
(351, 2, 'LIKE_MENU', 1, 'menus', '2025-11-22 06:44:19'),
(352, 13, 'USER_LOGIN', 13, 'users', '2025-11-22 06:45:00'),
(353, 13, 'UNLIKE_MENU', 1, 'menus', '2025-11-22 06:52:26'),
(354, 13, 'LIKE_MENU', 1, 'menus', '2025-11-22 06:52:28'),
(355, 13, 'BOOKMARK_MENU', 1, 'menus', '2025-11-22 06:52:41'),
(356, 13, 'UPDATE_PROFILE', 13, 'users', '2025-11-22 06:59:10'),
(357, 2, 'USER_LOGIN', 2, 'users', '2025-11-22 07:14:32'),
(358, 13, 'USER_LOGIN', 13, 'users', '2025-11-22 07:51:13'),
(359, 2, 'USER_LOGIN', 2, 'users', '2025-11-22 07:52:40'),
(360, 13, 'USER_LOGIN', 13, 'users', '2025-11-22 07:54:55'),
(361, 5, 'USER_LOGIN', 5, 'users', '2025-11-22 08:00:36'),
(362, 2, 'USER_LOGIN', 2, 'users', '2025-11-22 08:01:50'),
(363, 1, 'USER_LOGIN', 1, 'users', '2025-11-22 09:26:50'),
(364, 1, 'DELETE_MENU', 33, 'menus', '2025-11-22 09:27:09'),
(365, 2, 'USER_LOGIN', 2, 'users', '2025-11-22 09:27:41'),
(366, 2, 'UPDATE_MENU', 9, 'menus', '2025-11-22 10:55:30'),
(367, 2, 'UPDATE_MENU', 8, 'menus', '2025-11-22 10:56:15'),
(368, 2, 'UPDATE_MENU', 7, 'menus', '2025-11-22 10:57:19'),
(369, 2, 'UPDATE_MENU', 6, 'menus', '2025-11-22 10:58:22'),
(370, 2, 'UPDATE_MENU', 5, 'menus', '2025-11-22 10:59:17'),
(371, 2, 'UPDATE_MENU', 4, 'menus', '2025-11-22 11:00:11'),
(372, 2, 'UPDATE_MENU', 3, 'menus', '2025-11-22 11:00:51'),
(373, 13, 'USER_LOGIN', 13, 'users', '2025-11-22 11:45:17'),
(374, 2, 'USER_LOGIN', 2, 'users', '2025-11-22 11:46:34'),
(375, 2, 'LIKE_MENU', 2, 'menus', '2025-11-22 12:08:45'),
(376, 2, 'BOOKMARK_MENU', 1, 'menus', '2025-11-22 12:09:06'),
(377, 2, 'UNBOOKMARK_MENU', 1, 'menus', '2025-11-22 12:09:12'),
(378, 2, 'LIKE_MENU', 26, 'menus', '2025-11-22 12:09:16'),
(379, 5, 'USER_LOGIN', 5, 'users', '2025-11-22 12:49:33'),
(380, 5, 'REQUEST_PASSWORD_RESET', 5, 'users', '2025-11-22 12:50:38'),
(381, 1, 'USER_LOGIN', 1, 'users', '2025-11-22 12:51:47'),
(382, 1, 'LIKE_MENU', 1, 'menus', '2025-11-22 13:01:24'),
(383, 13, 'USER_LOGIN', 13, 'users', '2025-11-22 13:02:13'),
(384, 13, 'LIKE_MENU', 2, 'menus', '2025-11-22 13:02:49');

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
(1, 2, 'Ayam Bumbu Bawang', 'Resep ayam goreng sederhana dengan marinasi bumbu bawang yang gurih dan meresap.', '[\"1 dada ayam fillet, cuci bersih, potong kecil2\", \"1 buah jeruk nipis\", \"Bumbu marinate:\", \"5 siung bawang putih\", \"1 sdt lada hitam\", \"secukupnya Gula dan garam\", \"Minyam goreng\"]', '[\"Kucuri ayam dengan air jeruk nipis, kemudian marinate dengan bumbu halus selama 3 jam di kulkas\", \"Panaskan minyak, goreng ayam sampai matang dan kecoklatan\", \"Sajikan hangat dengan nasi putih.\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:09:44', '2025-11-09 16:10:24'),
(2, 2, 'Ayam Woku Manado', 'Resep Ayam Woku Manado yang lezat dan mudah dibuat.', '[\"1 Ekor Ayam Kampung (potong 12)--2 Buah Jeruk Nipis--2 Sdm Garam--3 Ruas Kunyit--7 Bawang Merah--7 Bawang Putih--10 Cabe Merah--10 Cabe Rawit Merah (sesuai selera)--3 Butir Kemiri--2 Batang Sereh--2 Lembar Daun Salam--2 Ikat Daun Kemangi--Penyedap Rasa--1 1/2 Gelas Air--\"]', '[\"Cuci bersih ayam dan tiriskan. Lalu peras jeruk nipis (kalo gak ada jeruk nipis bisa pake cuka) dan beri garam. Aduk hingga merata dan diamkan selama 5 menit\", \"biar ayam gak bau amis.--Goreng ayam tersebut setengah matang\", \"lalu tiriskan--Haluskan bumbu menggunakan blender. Bawang merah\", \"bawang putih\", \"cabe merah\", \"cabe rawit\", \"kemiri dan kunyit. Oh iya kasih minyak sedikit yaa biar bisa di blender. Untuk sereh nya di geprek aja terus di buat simpul.--Setelah bumbu di haluskan barulah di tumis. Jangan lupa sereh dan daun salamnya juga ikut di tumis. Di tumis sampai berubah warna ya 👌--Masukan ayam yang sudah di goreng setengah matang ke dalam bumbu yang sudah di tumis\", \"dan diamkan 5 menit dulu. Biar bumbu meresap. Lalu tuangkan 1 1/2 Gelas air. Lalu tambahkan penyedap rasa (saya 3 Sdt\", \"tapi sesuai selera ya) koreksi rasa dan Biar kan sampai mendidih--Setelah masakan mendidih\", \"lalu masukan daun kemangi yang sudah di potong potong. Masak lagi sekitar 10 menit. And taraaaaaaaaaaaaaa..... jadi deh Ayam Woku Manadonya.--Oh iyaa kalo mau di tambahkan potongan tomat merah juga bisa ko. Sesuai selera aja yaa buibuuuu 👌👌👌--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:17:49', '2025-11-09 16:10:13'),
(3, 2, 'Gurame Saus Padang', 'Resep Gurame Saus Padang yang lezat dan mudah dibuat.', '[\"Bahan utama:--1 ekor gurame--Bumbu untuk saus:--4 siung bawang putih (cincang halus)--3 siung bawang merah (cincang halus)--15 bh cabai merah (giling) banyaknya cabai sesuai selera ya😊--7 buah cabai rawit iris tipis (sesuai selera)--1/2 bawang bombang iris--Saus tiram--Saus tomat--Garam--Gula--Lada--1 bh wortel iris tipis memanjang (opsional)--1 buah tomat potong dadu (opsional)--1 batang irisan daun bawang--1 ruas jahe ukuran kecil iris tipis--Tepung maizena--250 ml air--Bumbu untuk menggoreng ikan :--Tepung beras--Tepung terigu--Garam--Jeruk nipis/ lemon--Minyak sayur--\"]', '[\"Cuci bersih ikan gurame yang akan dimasak. Setelah itu lumuri ikan dengan air perasan jeruk nipis/lemon diamkan beberapa menit untuk menghilangkan bau amisnya--Campur tepung terigu\", \"tepung beras dan garam menjadi satu\", \"kemudian lumuri ikan dengan tepung tersebut lalu goreng hingga matang lalu tiriskan. Pastikan minyak yang digunakan sudah panas ya biar ikannya tidak lengket😊.--Tumis bawang merah bawang putih hingga harum\", \"kemudian masukan cabai yang sudah giling\", \"aduk rata--Tambahkan bawang bombay\", \"saus tiram saus tomat dan jahe\", \"aduk sebentar lalu tambahkan 250ml air--Kemudian tambahkan garam dan gula sesuai selera.--Masukkan irisan cabai rawit dan daun bawang dan tomat. Aduk rata. Kemudian tambahkan larutan maizena agar saus mengental--Jika tidak ada tepung maizena. Diamkan beberapa menit hingga kandungan air menyusut dan saus terlihat mengental--Tes rasa. Jika sudah pas. Sajikan dan gurame saus padang siap disantap!😊--\"]', 'http://localhost:8080/assets/gurame-saus-padang-1.jpg', 'approved', '', '2025-10-23 15:18:12', '2025-11-22 11:00:51'),
(4, 2, 'Sate Kambing', 'Resep Sate Kambing yang lezat dan mudah dibuat.', '[\"Bahan-bahan :--500 gr Daging kambing--Daun pepaya--Bumbu yang dihaluskan :--4 siung bamer--4 siung baput--2 cm jahe--1 sdt ketumbar--1/2 sdt lada bubuk--1/2 sdt Garam--1 sdm air asam jawa--Bahan sambel kecap :--Kecap manis sesuai selera (aku pake bango)--Cabe rawit--Bamer--Tomat--Jeruk limau--\"]', '[\"1. Cuci bersih daging kambing\", \"potong\\\" kotak\", \"bungkus didalam lapisan daun pepaya (30menit)--2.setelah 30 menit\", \"keluarkan daging dr daun pepaya\", \"campurkan dengan bumbu yang sudah dihaluskan\", \"diamkan 15 menit--3. Susun daging pada tusuk sate kurlen 4-5 potong per tusuk--4. Bakar (bisa pakai arang / teflon)--5. Bolak balik sate hingga benar\\\" matang--6. Sajikan dengan sambal kecap--\"]', 'http://localhost:8080/assets/sate-kambing-1.jpg', 'approved', '', '2025-10-23 15:18:52', '2025-11-22 11:00:11'),
(5, 2, 'Beef Teriyaki', 'Resep Beef Teriyaki yang lezat dan mudah dibuat.', '[\"250 gr daging sapi--1 siung bawang bombai--5 siung bawang putih--1 sachet saus teriyaki (sy pk yg saor*, bsa diganti saus tiram)--1 sdm kecap manis--secukupnya garam--secukupnya lada--secukupnya gula--secukupnya penyedap rasa--\"]', '[\"Potong kecil-kecil memanjang daging sapi, lalu cuci bersih--Tambahkan garam dan lada diamkan selama kurleb 15 menit--Iris bawang bombai dan bawang putih--Tumis bawang bombai dan bawang putih, setelah harum masukkan daging--Tumis daging sebentar lalu tambahkan saus teriyaki, kecap manis, garam, lada, gula dan penyedap rasa, beri air sedikit--Tutup wajannya agar panasnya merata--Tunggu sampai daging matang, lalu garnish sesuai selera--Beef teriyaki siap untuk disantap--\"]', 'http://localhost:8080/assets/beef-teriyaki-1.jpg', 'approved', '', '2025-10-23 15:19:15', '2025-11-22 10:59:17'),
(6, 2, 'Martabak tahu pedas kulit lumpia', 'Resep Martabak tahu pedas kulit lumpia yang lezat dan mudah dibuat.', '[\"12 lembar kulit lumpia--4 buah tahu putih--2 batang daun bawang (iris halus)--1 batang daun seledri (iris halus)--3 btr telur (sy pakai 2 aja krn agak besar)--1 sdm terigu + 2 sdm air (sbgai lem kulit lumpia)--Sckpnya minyak utk menggoreng--Bumbu halus :--3 siung bawang putih--1/2 bks lada/merica bubuk--Sckpnya garam gula dan kaldu bubuk--3-4 cabe merah kriting iris serong--2 buah wortel agak kecil parut kasar--\"]', '[\"Haluskan tahu putih campur semua bahan kecuali kulit lumpia (icip rasa)--Ambil 1 lbr kulit lumpia beri 1 sdm adonan tahu ditengahnya lalu lipat (seperti amplop) dgn rapi, lalu beri lem disetiap sisisnya spy tdk terbuka pd saat digoreng, lakukan hingga habis--Panaskan minyak diatas wajan dgn api sedang cenderung kecil, goreng martabak sampai kecoklatan angkat dan tiriskan (abaikan wajannya yg sdh mulai jelek ya 😁)--Martabak tahu pedas kulit lumpia siap dihidangkan--Selamat mencoba\", \"😘😘--\"]', 'http://localhost:8080/assets/martabak-tahu-pedas-kulit-lumpia-1.jpg', 'approved', '', '2025-10-23 15:19:42', '2025-11-22 10:58:22'),
(7, 2, 'Orak arik telur buncis', 'Resep Orak arik telur buncis yang lezat dan mudah dibuat.', '[\"1 buah telur--10 buncis--3 bawang merah--2 bawang putih--6 cabe rawit--1 sdm kecap manis--Penyedap (garam gula atau penyedap lainnya)--\"]', '[\"Haluskan bawang merah\", \"bawang putih dan cabai--Iris buncis sesuai selera--Tumis bumbu halus hingga harum--Masukan buncis hingga agak layu--Masukan penyedap--Masukan telur\", \"lalu campur dengan buncis di orak arik hingga rata.. tunggu sebentar lalu di orak arik lagi.. supaya matangnya sempurna 👌(biar ga amis)--Masukan kecap manis--Masukan air sedikit biarkan menyusut dan meresap--Koreksi rasa--Jadi deh.. selamat menikmati 😇--\"]', 'http://localhost:8080/assets/orak-arik-telur-buncis-1.jpg', 'approved', '', '2025-10-23 15:20:07', '2025-11-22 10:57:19'),
(8, 2, 'Orek tempe manis pedas super simple', 'Resep Orek tempe manis pedas super simple yang lezat dan mudah dibuat.', '[\"Tempe (saya beli 3000 saja)--3 buah cabe gendot--3 buah cabe keriting--3 buah cabe rawit merah--2 siung bawang merah--1 siung bawang putih--Kecap--Saus tiram--Saus tomat--Kaldu bubuk--Garam--Gula--\"]', '[\"Potong2 tempe sesuai selera lalu goreng sebentar.--Tumis duo bawang dan semua cabe (cabe bisa sesuai selera ya pedesnya suka segimana) sampai harum. Masukkan tempe\", \"kecap\", \"saus tiram\", \"saus tomat\", \"kaldu bubuk\", \"gulgar\", \"dan tambahkan air satu gelas belimbing. Aduk2 dan tutup wajan.--Saat air sudah menyusut dan kuah kecap mengental. Cek rasa. Klo sudah oke. Matikan api. Sajikan. (Klo aku suka airnya sedikit sekali 😋). Selamat mencoba !--\"]', 'http://localhost:8080/assets/orek-tempe-manis-pedas-super-simple-1.jpg', 'approved', '', '2025-10-23 15:20:29', '2025-11-22 10:56:15'),
(9, 2, 'Lumpia udang kulit tahu', 'Resep Lumpia udang kulit tahu  yang lezat dan mudah dibuat.', '[\"50 gram ayam potong kotak kecil--200 gram udang besar potong kecil--1 lembar daun bawang iris halus--1/2 sendok teh garam halus/sesuai selera--1/2 sendok teh gula--1/4 sendok teh merica bubuk--1 sendok teh kecap ikan--1 sendok teh saus raja rasa--1 sendok teh saus tiram--1 sendok teh minyak wijen--1 sendok makan tepung sagu--1 sendok teh tepung terigu--1 butir telur ambil bagian putihnya saja, kocok lepas--2 lembar kembang tahu kering / kulit tahu direndam sebentar--\"]', '[\"Campur ayam & udang dengan semua bumbu & daun bawang\", \"aduk rata. Sesudah rata masukkan tepung sagu\", \"terigu & putih telur. Aduk rata.--Isikan ke lembaran kulit kembang tahu & lipat seperti melipat lumpia\", \"sambil dipadatkan agar rapi & bentuk agak pipih jangan terlalu bulat. Olesi putih telur diujung lipatan agar merekat rapat.--Kukus sekitar 15 menit hingga matang. Angkat dinginkan. Siap digoreng dengan api kecil sebentar. Anakku suka digoreng agak lama dengan kuning telur kocok sisa adonan tadi. Nanti jadi ada kayak jala2 crispynya yang gurih & kering. Angkat & tiriskan jika sudah digoreng matang.--Sajikan dengan dipotong diagonal (potong 2 serong). Siap disantap dengan saus sambal botol yang biasanya saya rebus dengan bawang goreng\", \"cabe rawit iris\", \"gula & sedikit air. Selamat mencoba.--\"]', 'http://localhost:8080/assets/lumpia-udang-kulit-tahu-1.jpg', 'approved', '', '2025-10-23 15:21:07', '2025-11-22 10:55:30'),
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
(23, 4, 'Telur Kornet', 'Resep Telur Kornet yang lezat dan mudah dibuat.', '[\"1/2 kaleng kornet--2-3 buah bawang prei--2 buah telur--1 sdm tepung terigu--secukupnya Cabe rawit--2 jumput Garam--1/2 sdt Merica--1/2 sdm Saus tiram--\"]', '[\"Campur kornet\", \"telur\", \"bawang prei dan cabe yg sudah di iris tipis. Tambahkan garam merica saos tiram secukupnya dan tepung terigu. Goreng 1 sdm dulu setelah matang baru cek rasa. Dirasa sudah pas. Adonan siap di goreng.😊--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 15:34:58', '2025-11-15 14:48:15'),
(24, 4, 'Penyet Tempe Sambel Korek Kemangi', 'Resep Penyet Tempe Sambel Korek Kemangi yang lezat dan mudah dibuat.', '[\"2 buah tempe--1 genggam Daun kemangi--15 cabe rawit--2 siung bawang putih--Gula--Garam--\"]', '[\"Iris tempe jd beberapa potong, bumbui dg bawang putih dan garam\", \"Goreng sampai kering--Ulek cabe rawit dan 1 siung bawang putih ukuran besar tambah gula garam sesuai selera--Setelah halus tuang minyak goreng panas bekas goreng tempe td--Penyet tempe diatas sambel dan beri daun kemangi--Tempe penyet siap disantap 😋😋😋--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'pending', '', '2025-10-23 15:35:11', '2025-10-23 15:35:11'),
(25, 4, 'Udang ala pop corn', 'Resep Udang ala pop corn yang lezat dan mudah dibuat.', '[\"1/4 kg udang basah ukuran sedang--1 bungkus kobe tepung ayam super crispy--secukupnya Air matang--Minyak untuk menggoreng--\"]', '[\"Buang kepala dan cangkang udang.--Cuci bersih udang.--Tepung bumbu dibagi jadi 2 adonan. Adonan basah dan adonan kering.--Masukkan udang ke dalam adonan tepung bumbu kering\", \"gulirkan ke dlm tepung sambil ditekan2. Lalu masukkan ke dalam Tepung bumbu adonan basah gulirkan lagi sambil di tekan2\", \"lalu masukkan lagi di adonan tepung bumbu kering Sambil ditekan2 agar tepung menempel sempurna.--Panaskan minyak goreng. Lalu goreng udang sampai kuning keemasan.--Angkat dan tiriskan.--Sajikan selagi hangat 😊--\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'rejected', 'Instruksi kurang jelas.', '2025-10-23 15:35:58', '2025-10-26 06:04:56'),
(26, 4, 'Ayam Geprek', 'Resep Ayam Geprek yang lezat dan mudah dibuat.', '[\"250 gr daging ayam (saya pakai fillet)\", \"Secukupnya gula dan garam\", \"50-100 gr tepung ayam serbaguna\", \"Secukupnya lalapan (kemangi,kol,timun)\", \"Secukupnya minyak panas\", \"❤sambal korek\", \"Secukupnya cabe rawit merah dan bwg putih\"]', '[\"Goreng ayam seperti ayam krispi\", \"Ulek semua bahan sambal kemudian campur dengan minyak panas bekas goreng ayam\", \"Geprek ayam kemudian campur dengan sambal,sajikan dengan lalapan ❤\"]', 'https://asset.kompas.com/crops/MZ_KjUJ4rxuZmCX1-_Kk3XplKyU=/32x0:1000x645/1200x800/data/photo/2022/02/11/62062e047e908.jpg', 'approved', '', '2025-10-23 16:23:16', '2025-10-26 06:02:33'),
(31, 2, 'semur daging sapi', 'semur daging sapi enak pisan', '[\"1kg daging\"]', '[\"rebus daging\"]', 'http://localhost:8080/assets/semur-daging-sapi-1.jpg', 'rejected', 'resep tidak jelas, bahan bahan dan langkah langkah tidak jelas', '2025-11-09 14:54:22', '2025-11-15 15:01:33');

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
(1, 1),
(2, 1),
(26, 1),
(3, 2),
(4, 3),
(5, 4),
(31, 4),
(6, 5),
(9, 5),
(7, 6),
(8, 7),
(9, 8),
(9, 20),
(26, 20),
(31, 22),
(31, 40),
(6, 50),
(26, 50);

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
(14, 5, 2, 0, 0, '2025-11-08 09:54:39', '2025-11-09 11:59:11'),
(20, 11, 2, 1, 0, '2025-11-08 10:10:58', '2025-11-09 07:48:35'),
(21, 10, 2, 1, 0, '2025-11-08 10:11:08', '2025-11-09 07:17:49'),
(22, 1, 2, 1, 0, '2025-11-09 07:48:41', '2025-11-22 06:44:19'),
(23, 7, 2, 1, 0, '2025-11-09 07:59:54', '2025-11-09 10:28:59'),
(24, 6, 2, 1, 0, '2025-11-09 11:08:57', '2025-11-09 11:08:57'),
(25, 15, 2, 0, 0, '2025-11-09 11:15:57', '2025-11-22 06:29:49'),
(26, 4, 2, 0, 0, '2025-11-09 12:07:25', '2025-11-22 06:26:34'),
(27, 12, 2, 1, 0, '2025-11-09 12:07:29', '2025-11-09 12:11:41'),
(28, 1, 13, 1, 0, '2025-11-21 10:38:40', '2025-11-22 06:52:28'),
(29, 15, 13, 1, 0, '2025-11-21 10:39:48', '2025-11-21 10:39:48'),
(30, 2, 2, 1, 0, '2025-11-22 12:08:45', '2025-11-22 12:08:45'),
(31, 26, 2, 1, 0, '2025-11-22 12:09:15', '2025-11-22 12:09:15'),
(32, 1, 1, 1, 0, '2025-11-22 13:01:24', '2025-11-22 13:01:24'),
(33, 2, 13, 1, 0, '2025-11-22 13:02:49', '2025-11-22 13:02:49');

-- --------------------------------------------------------

--
-- Table structure for table `notifications`
--

CREATE TABLE `notifications` (
  `notification_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `message` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` enum('info','success','warning','danger') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'info',
  `is_read` tinyint(1) NOT NULL DEFAULT '0',
  `related_id` bigint UNSIGNED DEFAULT NULL,
  `related_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `notifications`
--

INSERT INTO `notifications` (`notification_id`, `user_id`, `title`, `message`, `type`, `is_read`, `related_id`, `related_type`, `created_at`) VALUES
(1, 4, 'Resep Ditolak', 'Resep \"Udang ala pop corn\" ditolak. Alasan: Instruksi kurang jelas.', 'danger', 1, 25, 'menu', '2025-10-26 06:04:56'),
(2, 2, 'Resep Ditolak', 'Resep \"semur daging sapi\" ditolak. Alasan: resep tidak jelas, bahan bahan dan langkah langkah tidak jelas', 'danger', 1, 31, 'menu', '2025-11-15 15:01:33'),
(4, 2, 'Resep Ditolak', 'Resep \"Ayam Geprek\" ditolak. Alasan: Bahan bahan dan intruksi tidak jelas tolong perbaiki yaa', 'danger', 1, 33, 'menu', '2025-11-21 14:38:05'),
(5, 2, 'Komentar Baru', 'raka wilangga berkomentar: \"mantap mudah di pahami resep nya\"', 'info', 1, 1, 'menu', '2025-11-22 07:14:10'),
(6, 13, 'Balasan Baru', 'Budi Santoso membalas komentar Anda: \"terima kasih banyak semoga bermanfaat\"', 'info', 1, 1, 'menu', '2025-11-22 07:53:20');

-- --------------------------------------------------------

--
-- Table structure for table `password_resets`
--

CREATE TABLE `password_resets` (
  `reset_id` int NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `verification_code` varchar(6) COLLATE utf8mb4_unicode_ci NOT NULL,
  `expires_at` datetime NOT NULL,
  `is_used` tinyint(1) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `password_resets`
--

INSERT INTO `password_resets` (`reset_id`, `user_id`, `verification_code`, `expires_at`, `is_used`, `created_at`) VALUES
(9, 1, '053082', '2025-11-17 10:29:00', 1, '2025-11-17 10:19:00'),
(10, 1, '835542', '2025-11-21 16:40:29', 0, '2025-11-21 16:30:29'),
(11, 13, '031214', '2025-11-21 17:42:14', 1, '2025-11-21 17:32:14'),
(12, 5, '041671', '2025-11-22 20:00:36', 0, '2025-11-22 19:50:36');

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
(1, 1, 'active', 'Miqdam syiam nurrohmann', 'damnihbos4@gmail.com', '$2a$10$/5mtuudIOhpzj9Gch/ZVnu/2SB2THRYsW9bXPQqfqMGxCvRbDkrV.', '2025-10-23 13:55:32', '2025-11-17 03:19:15'),
(2, 3, 'active', 'Budi Santoso', 'budi@gmail.com', '$2a$10$4tSbPmmnpscN2Gfyse9GieU8aTPeQMrQC/EqxCSTG/urVFANwsMCq', '2025-10-23 13:56:42', '2025-11-22 12:09:06'),
(3, 3, 'active', 'Ani Lestari', 'ani@gmail.com', '$2a$10$c5pLKPfMoO2h/5mYVq/wieFfeseTHThLFz7W9F0j3EQi4SqOrmMUW', '2025-10-23 13:57:06', '2025-10-26 07:10:07'),
(4, 3, 'active', 'Candra Wijaya', 'candra@gmail.com', '$2a$10$LcnkpHDgymYKSPI2Ri5qoetb2l4bRMrppSArtF/q6GORwTY7TmzYO', '2025-10-23 13:57:25', '2025-11-21 14:15:44'),
(5, 2, 'active', 'Nizar Akmal', 'nizarakmal8@gmail.com', '$2a$10$FX8E3Pplm58noZ44FLGH.ePT5ecVcjZmypNRrY3rxtbVmVQSeFGc.', '2025-10-23 13:58:08', '2025-11-16 08:43:52'),
(6, 3, 'active', 'Asep Wiyanto', 'asep@gmail.com', '$2a$10$p1/AQxEl50ml2sAgl4YdperlBm98JXH1pWijGt.r8P6cblDWKlXJe', '2025-10-26 06:25:26', '2025-11-16 09:11:18'),
(8, 3, 'active', 'Akbar wijaya', 'akbar@gmail.com', '$2a$10$joaWEutecWGhWKUtq9PlMuNZVu5UHoj6zmOt.bk1cE6AYDEM2oTNC', '2025-11-16 09:14:34', '2025-11-16 09:20:16'),
(9, 3, 'active', 'Hadi Yusuf', 'yusuf@gmail.com', '$2a$10$VRGVzWSwB2ESpdoa3woCoeToJNGr.4TM.7.g.fPxX4U6t8239Gmwi', '2025-11-16 09:22:04', '2025-11-16 09:22:04'),
(10, 3, 'inactive', 'Rifky Najra', 'rifky@gmail.com', '$2a$10$o2h9jA5hCJXZKioPoqBbneuvfq8WOQYobLR.Qpt0roWBILXvHlB1u', '2025-11-16 09:22:33', '2025-11-16 09:27:22'),
(11, 3, 'inactive', 'Muflih Afif', 'afif@gmail.com', '$2a$10$noMKKZsmNbMNsP1aXwzUXuL92mPmKGKgu5MV0LY2lDovHABUrHW9G', '2025-11-16 09:23:03', '2025-11-16 09:27:18'),
(12, 3, 'active', 'sam', 'sam@gmail.com', '$2a$10$zjKc3ZjgxV.HEgHD9ndRxOl2awFb4wguIadzHgLfb87eSU5bmKoN6', '2025-11-16 09:27:06', '2025-11-16 09:27:06'),
(13, 3, 'active', 'raka wilangga', 'rakawilangga713@gmail.com', '$2a$10$Fg7Qn0IHcOF2jabmw8DlpegXRvgQV7sPl9N5uGYFAE6isj9uJmU8e', '2025-11-21 10:27:31', '2025-11-22 06:59:10');

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
(13, 1),
(3, 2),
(4, 2),
(6, 2),
(2, 3),
(4, 3),
(6, 3),
(2, 7),
(2, 8),
(2, 11),
(2, 12);

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
(1, 1, 'admin ganteng pisann', '/assets/profiles/profile_1_1763349504.png'),
(2, 2, '', ''),
(3, 3, '', ''),
(4, 4, '', ''),
(5, 5, '', ''),
(6, 6, '', ''),
(8, 8, '', ''),
(9, 9, '', ''),
(10, 10, '', ''),
(11, 11, '', ''),
(12, 12, '', ''),
(13, 13, 'raka ', '/assets/profiles/profile_13_1763721096.png');

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
  ADD KEY `fk_comments_user` (`user_id`),
  ADD KEY `idx_comments_parent_id` (`parent_id`),
  ADD KEY `idx_comments_menu_parent` (`menu_id`,`parent_id`);

--
-- Indexes for table `email_verifications`
--
ALTER TABLE `email_verifications`
  ADD PRIMARY KEY (`verification_id`),
  ADD KEY `idx_user_id` (`user_id`),
  ADD KEY `idx_verification_code` (`verification_code`),
  ADD KEY `idx_is_used` (`is_used`);

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
-- Indexes for table `notifications`
--
ALTER TABLE `notifications`
  ADD PRIMARY KEY (`notification_id`),
  ADD KEY `idx_user_id` (`user_id`),
  ADD KEY `idx_is_read` (`is_read`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `idx_user_read` (`user_id`,`is_read`);

--
-- Indexes for table `password_resets`
--
ALTER TABLE `password_resets`
  ADD PRIMARY KEY (`reset_id`),
  ADD KEY `idx_user_id` (`user_id`),
  ADD KEY `idx_verification_code` (`verification_code`),
  ADD KEY `idx_is_used` (`is_used`);

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
  MODIFY `category_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `comments`
--
ALTER TABLE `comments`
  MODIFY `comment_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `email_verifications`
--
ALTER TABLE `email_verifications`
  MODIFY `verification_id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `log_activity`
--
ALTER TABLE `log_activity`
  MODIFY `activity_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=385;

--
-- AUTO_INCREMENT for table `menus`
--
ALTER TABLE `menus`
  MODIFY `menu_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=34;

--
-- AUTO_INCREMENT for table `menu_votes`
--
ALTER TABLE `menu_votes`
  MODIFY `vote_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=34;

--
-- AUTO_INCREMENT for table `notifications`
--
ALTER TABLE `notifications`
  MODIFY `notification_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `password_resets`
--
ALTER TABLE `password_resets`
  MODIFY `reset_id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `roles`
--
ALTER TABLE `roles`
  MODIFY `role_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `tags`
--
ALTER TABLE `tags`
  MODIFY `tag_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=57;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `user_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=14;

--
-- AUTO_INCREMENT for table `user_profiles`
--
ALTER TABLE `user_profiles`
  MODIFY `profile_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=14;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `comments`
--
ALTER TABLE `comments`
  ADD CONSTRAINT `fk_comments_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_comments_parent` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`comment_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

--
-- Constraints for table `email_verifications`
--
ALTER TABLE `email_verifications`
  ADD CONSTRAINT `fk_email_verifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

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
-- Constraints for table `notifications`
--
ALTER TABLE `notifications`
  ADD CONSTRAINT `fk_notification_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

--
-- Constraints for table `password_resets`
--
ALTER TABLE `password_resets`
  ADD CONSTRAINT `fk_password_resets_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE;

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
