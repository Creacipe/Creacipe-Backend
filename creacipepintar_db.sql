-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Oct 20, 2025 at 09:47 AM
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
-- Table structure for table `menus`
--

CREATE TABLE `menus` (
  `menu_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text,
  `ingredients` json DEFAULT NULL,
  `instructions` text,
  `image_url` varchar(255) DEFAULT NULL,
  `status` enum('pending','approved','rejected') NOT NULL DEFAULT 'pending',
  `rejection_reason` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `menus`
--

INSERT INTO `menus` (`menu_id`, `user_id`, `title`, `description`, `ingredients`, `instructions`, `image_url`, `status`, `rejection_reason`, `created_at`, `updated_at`) VALUES
(101, 10, 'Rendang Sapi', 'Daging sapi empuk dengan bumbu rempah khas Padang yang medok.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(102, 10, 'Ayam Bakar Madu', 'Ayam bakar dengan olesan madu dan kecap yang manis dan gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(103, 11, 'Sayur Asem Jakarta', 'Sayur asem dengan kuah segar berisi berbagai macam sayuran.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(104, 11, 'Ikan Goreng Tepung', 'Ikan dori fillet yang digoreng renyah dengan balutan tepung.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(105, 12, 'Nasi Goreng Seafood', 'Nasi goreng dengan isian udang, cumi, dan bumbu spesial.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(106, 10, 'Mie Goreng Jawa', 'Mie kuning yang dimasak dengan bumbu kemiri dan sayuran.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(107, 11, 'Tahu Gejrot', 'Tahu pong diguyur kuah asam pedas dari bawang dan cabai.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(108, 12, 'Tempe Bacem', 'Tempe yang direbus dengan air kelapa dan gula merah, lalu digoreng.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(109, 10, 'Sate Ayam Madura', 'Sate ayam dengan bumbu kacang khas Madura yang kental.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(110, 11, 'Bubur Sumsum', 'Bubur lembut dari tepung beras disajikan dengan saus gula merah.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(111, 10, 'Soto Ayam Lamongan', 'Soto bening dengan suwiran ayam, koya gurih, dan jeruk nipis.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(112, 11, 'Rawon Sapi Surabaya', 'Sup hitam khas Jawa Timur dengan kluwek dan daging sapi empuk.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(113, 12, 'Gado-Gado Bumbu Kacang', 'Sayuran rebus disiram bumbu kacang kental, lengkap dengan kerupuk.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(114, 10, 'Ikan Bakar Rica', 'Ikan segar dibakar lalu dilapisi sambal rica pedas menyengat.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(115, 11, 'Nasi Uduk Betawi', 'Nasi gurih santan dengan lauk pelengkap dan sambal.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(116, 12, 'Mie Rebus Kampung', 'Mie telur direbus dengan kuah kaldu gurih dan sayuran.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(117, 10, 'Tahu Isi Sayur', 'Tahu goreng berisi irisan sayur berbumbu lalu disajikan dengan cabai.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(118, 11, 'Tempe Mendoan', 'Tempe tipis berbalut tepung digoreng setengah matang, lembut dan gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(119, 12, 'Bakso Sapi Kuah', 'Bakso kenyal dalam kuah kaldu hangat dengan mie dan sayuran.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(120, 10, 'Ayam Goreng Lengkuas', 'Ayam goreng renyah dengan taburan lengkuas garing.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(121, 11, 'Sapi Lada Hitam', 'Irisan daging sapi ditumis dengan saus lada hitam yang gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(122, 12, 'Ikan Kuah Kuning', 'Ikan dimasak kuah kunyit segar dengan rempah sederhana.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(123, 10, 'Capcay Tumis', 'Aneka sayuran ditumis cepat dengan bawang dan saus ringan.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(124, 11, 'Nasi Liwet Sunda', 'Nasi gurih dengan teri, cabai, dan daun salam khas Sunda.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(125, 12, 'Mie Goreng Seafood Pedas', 'Mie goreng dengan udang dan cumi, pedas dan wangi.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(126, 10, 'Tahu Bacem Goreng', 'Tahu yang dibacem manis gurih lalu digoreng hingga kecokelatan.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(127, 11, 'Tempe Orek Kecap', 'Tempe iris tipis dimasak manis pedas dengan kecap.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(128, 12, 'Ayam Woku', 'Ayam berbumbu woku pedas segar khas Manado.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(129, 10, 'Sate Sapi Kecap', 'Sate daging sapi empuk dibakar dengan olesan kecap manis.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(130, 11, 'Ikan Asam Pedas', 'Ikan dimasak kuah asam pedas segar menggugah selera.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(131, 12, 'Sayur Lodeh', 'Sayur santan lembut dengan labu, kacang panjang, dan tempe.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(132, 10, 'Nasi Kebuli Sapi', 'Nasi berbumbu rempah kuat dengan potongan daging sapi.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(133, 11, 'Mie Tek-Tek Goreng', 'Mie goreng kaki lima dengan telur, sayur, dan bumbu tajam.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(134, 12, 'Perkedel Kentang', 'Perkedel kentang lembut, gurih, nikmat untuk lauk samping.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(135, 10, 'Ayam Penyet Sambal Terasi', 'Ayam goreng dipenyet disajikan dengan sambal terasi pedas.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(136, 11, 'Sapi Semur Betawi', 'Semur daging sapi manis gurih dengan rempah hangat.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(137, 12, 'Ikan Pepes Daun Pisang', 'Ikan dibumbui lalu dikukus/pepes dalam daun pisang wangi.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(138, 10, 'Cah Kangkung Bawang Putih', 'Kangkung tumis cepat dengan bawang putih, sederhana tapi nikmat.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(139, 11, 'Nasi Goreng Kampung', 'Nasi goreng sederhana dengan teri, cabai, dan kemangi.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(140, 12, 'Mie Aceh Tumis', 'Mie kuning bumbu rempah kuat dengan irisan daging.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(141, 10, 'Tahu Telur Surabaya', 'Tahu telur dengan saus petis manis gurih dan sayuran.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(142, 11, 'Tempe Garit Krispi', 'Tempe diserut tipis digoreng kriuk sebagai camilan atau lauk.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(143, 12, 'Ayam Taliwang', 'Ayam bakar pedas khas Lombok, rasa kuat dan smoky.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(144, 10, 'Sapi Balado', 'Irisan daging sapi digoreng lalu dibalut sambal balado.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(145, 11, 'Ikan Goreng Sambal Matah', 'Ikan goreng garing disajikan dengan sambal matah segar.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(146, 12, 'Urap Sayur', 'Sayuran kukus dengan kelapa parut berbumbu wangi.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(147, 10, 'Nasi Kuning Komplit', 'Nasi kunyit gurih dengan lauk telur, ayam, dan sambal.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(148, 11, 'Mie Kuah Bakso', 'Mie kuah hangat dengan bakso dan sawi, cocok saat hujan.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(149, 12, 'Bakwan Tahu Sayur', 'Gorengan renyah berisi tahu dan parutan sayuran berbumbu.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(150, 10, 'Ayam Rica-Rica', 'Ayam pedas dengan bumbu rica merah segar.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(151, 11, 'Sapi Tumis Jamur', 'Sapi ditumis cepat dengan jamur dan bawang bombai.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(152, 12, 'Ikan Bakar Jimbaran', 'Ikan bakar oles sambal khas Jimbaran, manis pedas.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(153, 10, 'Sayur Bayam Bening', 'Sayur bening bayam dan jagung, segar menenangkan.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(154, 11, 'Nasi Goreng Kambing (Sapi Varian)', 'Versi sapi dari nasi goreng bumbu kebuli yang gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(155, 12, 'Mie Goreng Tekstur Rumahan', 'Mie goreng harian simpel, gurih dan wangi.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(156, 10, 'Tahu Sambal Kecap', 'Tahu goreng disiram sambal kecap pedas manis.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(157, 11, 'Tempe Bacem Iris Tipis', 'Tempe bacem manis gurih, dipotong tipis agar meresap.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(158, 12, 'Ayam Gulai Santan', 'Ayam dimasak santan berbumbu gulai hangat.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(159, 10, 'Sapi Tumis Cabe Ijo', 'Sapi iris tipis ditumis dengan cabai hijau pedas.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(160, 11, 'Ikan Pindang Serani', 'Ikan dimasak kuah asam segar dengan tomat dan sereh.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(161, 12, 'Tumis Toge Tahu', 'Tauge dan tahu ditumis cepat, sederhana dan sehat.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(162, 10, 'Nasi Tutug Oncom', 'Nasi diaduk oncom panggang, wangi khas Sunda.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(163, 11, 'Mie Ayam Jamur', 'Mie ayam gurih topping jamur kecap.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(164, 12, 'Tahu Lada Garam', 'Tahu renyah dilumuri bawang putih, garam, dan cabai.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(165, 10, 'Ayam Kecap Pedas', 'Ayam dimasak kecap manis pedas, cocok untuk makan malam.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(166, 11, 'Sapi Teriyaki Lokal', 'Sapi tumis manis gurih gaya rumahan.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(167, 12, 'Ikan Kukus Jahe', 'Ikan dikukus dengan jahe dan bawang, rasa ringan.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(168, 10, 'Oseng Buncis Wortel', 'Buncis dan wortel oseng cepat dengan bawang.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(169, 11, 'Nasi Goreng Petai', 'Nasi goreng dengan petai dan irisan cabai, wangi kuat.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(170, 12, 'Mie Goreng Pedas Manis', 'Mie goreng dengan komposisi pedas manis seimbang.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(171, 10, 'Tahu Cabe Garam', 'Tahu goreng kriuk dengan cabai dan garam, camilan favorit.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(172, 11, 'Tempe Penyet Sambal Ijo', 'Tempe goreng dipenyet sambal hijau pedas.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(173, 12, 'Ayam Bakar Kecap', 'Ayam dibakar oles kecap manis, smoky dan legit.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(174, 10, 'Sapi Geprek Bawang', 'Daging sapi digeprek dengan bawang pedas gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(175, 11, 'Ikan Sambal Dabu-Dabu', 'Ikan goreng disajikan dengan sambal dabu-dabu segar.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(176, 12, 'Sayur Sop Rumahan', 'Sop sayur bening dengan wortel, kentang, dan kol.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(177, 10, 'Nasi Goreng Teri Medan', 'Nasi goreng dengan teri medan, renyah gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(178, 11, 'Mie Godog Jawa', 'Mie rebus kuah Jawa dengan telur dan sayuran.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(179, 12, 'Tahu Kuah Kecap', 'Tahu direbus dalam kuah kecap manis gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(180, 10, 'Ayam Kalasan', 'Ayam berbumbu manis gurih khas Yogyakarta.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(181, 11, 'Sapi Oseng Mercon', 'Oseng daging sapi super pedas bagi pencinta cabai.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(182, 12, 'Ikan Woku Belanga', 'Ikan masak woku berkuah, pedas segar.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(183, 10, 'Tumis Sawi Putih Telur', 'Sawi putih ditumis dengan telur orak-arik, simple dan enak.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(184, 11, 'Nasi Goreng Telur Asin', 'Nasi goreng dengan saus telur asin gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(185, 12, 'Mie Nyemek', 'Mie kuah kental khas Jawa, pedas gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(186, 10, 'Tahu Kecap Bawang', 'Tahu goreng disiram kecap bawang sederhana.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(187, 11, 'Tempe Kering Balado', 'Tempe kering kriuk dibalut saus balado pedas manis.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(188, 12, 'Ayam Opor Putih', 'Opor ayam santan putih lembut dan hangat.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(189, 10, 'Sapi Bistik Rumahan', 'Bistik daging sapi gaya rumahan manis gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(190, 11, 'Ikan Goreng Sambal Terasi', 'Ikan goreng garing dengan sambal terasi pedas.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(191, 12, 'Urap Kencur', 'Sayur urap dengan aroma kencur lebih kuat.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(192, 10, 'Nasi Jagung Komplit', 'Nasi jagung dengan lauk sederhana khas desa.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(193, 11, 'Mie Goreng Baso Sosis', 'Mie goreng dengan bakso dan sosis untuk porsi kenyang.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(194, 12, 'Tahu Pong Kuah', 'Tahu pong disajikan dalam kuah segar asam gurih.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(195, 10, 'Ayam Cabe Ijo', 'Ayam dimasak sambal cabai hijau pedas.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(196, 11, 'Sapi Tumis Paprika', 'Daging sapi tumis paprika renyah dan manis alami.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(197, 12, 'Ikan Kuah Asam Belimbing', 'Ikan dimasak kuah asam belimbing wuluh segar.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(198, 10, 'Sayur Asem Sunda', 'Sayur asem bening asam segar khas Sunda.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(199, 11, 'Nasi Goreng Mawut', 'Campuran nasi dan mie digoreng bersama jadi satu.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL),
(200, 12, 'Mie Goreng Jawa Pedas', 'Mie goreng gaya Jawa dengan tingkat pedas menantang.', NULL, NULL, NULL, 'approved', NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `menu_ratings`
--

CREATE TABLE `menu_ratings` (
  `rating_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `menu_id` bigint UNSIGNED NOT NULL,
  `rating` tinyint UNSIGNED NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `menu_ratings`
--

INSERT INTO `menu_ratings` (`rating_id`, `user_id`, `menu_id`, `rating`, `created_at`, `updated_at`) VALUES
(1, 10, 101, 5, NULL, NULL),
(2, 10, 109, 4, NULL, NULL),
(3, 10, 111, 4, NULL, NULL),
(4, 10, 112, 5, NULL, NULL),
(5, 10, 119, 4, NULL, NULL),
(6, 10, 120, 5, NULL, NULL),
(7, 10, 121, 4, NULL, NULL),
(8, 10, 128, 4, NULL, NULL),
(9, 10, 129, 4, NULL, NULL),
(10, 10, 135, 5, NULL, NULL),
(11, 10, 136, 4, NULL, NULL),
(12, 10, 143, 4, NULL, NULL),
(13, 10, 144, 5, NULL, NULL),
(14, 10, 158, 5, NULL, NULL),
(15, 10, 159, 4, NULL, NULL),
(16, 10, 173, 5, NULL, NULL),
(17, 10, 174, 5, NULL, NULL),
(18, 10, 180, 4, NULL, NULL),
(19, 10, 181, 4, NULL, NULL),
(20, 10, 188, 5, NULL, NULL),
(21, 10, 195, 5, NULL, NULL),
(22, 10, 196, 5, NULL, NULL),
(23, 11, 101, 4, NULL, NULL),
(24, 11, 108, 4, NULL, NULL),
(25, 11, 111, 5, NULL, NULL),
(26, 11, 112, 4, NULL, NULL),
(27, 11, 115, 4, NULL, NULL),
(28, 11, 116, 5, NULL, NULL),
(29, 11, 119, 4, NULL, NULL),
(30, 11, 124, 5, NULL, NULL),
(31, 11, 130, 5, NULL, NULL),
(32, 11, 131, 5, NULL, NULL),
(33, 11, 132, 5, NULL, NULL),
(34, 11, 137, 5, NULL, NULL),
(35, 11, 146, 4, NULL, NULL),
(36, 11, 147, 5, NULL, NULL),
(37, 11, 148, 4, NULL, NULL),
(38, 11, 167, 5, NULL, NULL),
(39, 11, 178, 5, NULL, NULL),
(40, 11, 179, 5, NULL, NULL),
(41, 11, 182, 4, NULL, NULL),
(42, 11, 185, 5, NULL, NULL),
(43, 11, 188, 4, NULL, NULL),
(44, 11, 194, 5, NULL, NULL),
(45, 11, 197, 5, NULL, NULL),
(46, 11, 198, 4, NULL, NULL),
(47, 12, 101, 4, NULL, NULL),
(48, 12, 102, 5, NULL, NULL),
(49, 12, 109, 5, NULL, NULL),
(50, 12, 111, 4, NULL, NULL),
(51, 12, 112, 5, NULL, NULL),
(52, 12, 121, 5, NULL, NULL),
(53, 12, 129, 5, NULL, NULL),
(54, 12, 136, 5, NULL, NULL),
(55, 12, 143, 5, NULL, NULL),
(56, 12, 150, 5, NULL, NULL),
(57, 12, 151, 5, NULL, NULL),
(58, 12, 158, 5, NULL, NULL),
(59, 12, 159, 4, NULL, NULL),
(60, 12, 173, 5, NULL, NULL),
(61, 12, 174, 5, NULL, NULL),
(62, 12, 195, 4, NULL, NULL),
(63, 12, 196, 5, NULL, NULL),
(64, 13, 102, 5, '2025-10-18 16:37:44.341', '2025-10-18 16:37:44.341');

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
(102, 1),
(109, 1),
(111, 1),
(120, 1),
(128, 1),
(135, 1),
(143, 1),
(150, 1),
(158, 1),
(165, 1),
(173, 1),
(180, 1),
(188, 1),
(195, 1),
(101, 2),
(112, 2),
(119, 2),
(121, 2),
(129, 2),
(136, 2),
(144, 2),
(151, 2),
(159, 2),
(166, 2),
(174, 2),
(181, 2),
(189, 2),
(196, 2),
(104, 3),
(114, 3),
(122, 3),
(130, 3),
(137, 3),
(145, 3),
(152, 3),
(160, 3),
(167, 3),
(175, 3),
(182, 3),
(190, 3),
(197, 3),
(103, 4),
(113, 4),
(123, 4),
(131, 4),
(134, 4),
(138, 4),
(146, 4),
(153, 4),
(161, 4),
(168, 4),
(176, 4),
(183, 4),
(191, 4),
(198, 4),
(105, 5),
(115, 5),
(124, 5),
(132, 5),
(139, 5),
(147, 5),
(154, 5),
(162, 5),
(169, 5),
(177, 5),
(184, 5),
(192, 5),
(199, 5),
(106, 6),
(116, 6),
(125, 6),
(133, 6),
(140, 6),
(148, 6),
(155, 6),
(163, 6),
(170, 6),
(178, 6),
(185, 6),
(193, 6),
(199, 6),
(200, 6),
(107, 7),
(117, 7),
(126, 7),
(141, 7),
(149, 7),
(156, 7),
(161, 7),
(164, 7),
(171, 7),
(179, 7),
(186, 7),
(194, 7),
(108, 8),
(118, 8),
(127, 8),
(142, 8),
(157, 8),
(172, 8),
(187, 8),
(104, 10),
(105, 10),
(117, 10),
(118, 10),
(120, 10),
(125, 10),
(126, 10),
(133, 10),
(134, 10),
(135, 10),
(139, 10),
(141, 10),
(142, 10),
(145, 10),
(149, 10),
(154, 10),
(155, 10),
(164, 10),
(169, 10),
(170, 10),
(171, 10),
(172, 10),
(174, 10),
(175, 10),
(177, 10),
(184, 10),
(187, 10),
(189, 10),
(190, 10),
(193, 10),
(199, 10),
(200, 10),
(102, 11),
(109, 11),
(114, 11),
(129, 11),
(143, 11),
(152, 11),
(173, 11),
(103, 12),
(108, 12),
(113, 12),
(115, 12),
(116, 12),
(124, 12),
(126, 12),
(132, 12),
(136, 12),
(137, 12),
(146, 12),
(147, 12),
(157, 12),
(167, 12),
(178, 12),
(180, 12),
(185, 12),
(191, 12),
(192, 12),
(106, 13),
(121, 13),
(123, 13),
(127, 13),
(128, 13),
(138, 13),
(140, 13),
(144, 13),
(150, 13),
(151, 13),
(156, 13),
(159, 13),
(161, 13),
(162, 13),
(163, 13),
(165, 13),
(166, 13),
(168, 13),
(181, 13),
(183, 13),
(186, 13),
(195, 13),
(196, 13),
(101, 14),
(103, 14),
(111, 14),
(112, 14),
(116, 14),
(119, 14),
(122, 14),
(130, 14),
(131, 14),
(148, 14),
(153, 14),
(158, 14),
(160, 14),
(176, 14),
(178, 14),
(179, 14),
(182, 14),
(185, 14),
(188, 14),
(194, 14),
(197, 14),
(198, 14),
(107, 20),
(114, 20),
(125, 20),
(128, 20),
(130, 20),
(135, 20),
(139, 20),
(140, 20),
(143, 20),
(144, 20),
(145, 20),
(150, 20),
(159, 20),
(164, 20),
(165, 20),
(169, 20),
(170, 20),
(171, 20),
(172, 20),
(174, 20),
(175, 20),
(181, 20),
(182, 20),
(185, 20),
(187, 20),
(190, 20),
(195, 20),
(200, 20),
(102, 21),
(106, 21),
(108, 21),
(110, 21),
(126, 21),
(127, 21),
(129, 21),
(136, 21),
(141, 21),
(147, 21),
(152, 21),
(156, 21),
(157, 21),
(163, 21),
(165, 21),
(166, 21),
(170, 21),
(173, 21),
(179, 21),
(180, 21),
(186, 21),
(187, 21),
(188, 21),
(189, 21),
(101, 22),
(102, 22),
(104, 22),
(105, 22),
(106, 22),
(107, 22),
(109, 22),
(111, 22),
(112, 22),
(114, 22),
(115, 22),
(117, 22),
(118, 22),
(119, 22),
(120, 22),
(121, 22),
(123, 22),
(124, 22),
(125, 22),
(127, 22),
(128, 22),
(129, 22),
(131, 22),
(132, 22),
(133, 22),
(134, 22),
(136, 22),
(137, 22),
(138, 22),
(139, 22),
(140, 22),
(141, 22),
(142, 22),
(143, 22),
(144, 22),
(145, 22),
(147, 22),
(148, 22),
(149, 22),
(150, 22),
(151, 22),
(152, 22),
(154, 22),
(155, 22),
(158, 22),
(159, 22),
(162, 22),
(163, 22),
(164, 22),
(165, 22),
(166, 22),
(168, 22),
(169, 22),
(170, 22),
(171, 22),
(173, 22),
(174, 22),
(175, 22),
(176, 22),
(177, 22),
(178, 22),
(180, 22),
(181, 22),
(182, 22),
(183, 22),
(184, 22),
(185, 22),
(186, 22),
(187, 22),
(188, 22),
(189, 22),
(190, 22),
(192, 22),
(193, 22),
(195, 22),
(196, 22),
(199, 22),
(200, 22),
(103, 23),
(107, 23),
(113, 23),
(122, 23),
(130, 23),
(146, 23),
(153, 23),
(160, 23),
(167, 23),
(191, 23),
(194, 23),
(197, 23),
(198, 23),
(105, 30),
(115, 30),
(117, 30),
(142, 30),
(147, 30),
(149, 30),
(153, 30),
(164, 30),
(171, 30),
(187, 30),
(192, 30),
(103, 31),
(104, 31),
(111, 31),
(113, 31),
(116, 31),
(118, 31),
(119, 31),
(122, 31),
(123, 31),
(130, 31),
(134, 31),
(137, 31),
(138, 31),
(146, 31),
(148, 31),
(151, 31),
(157, 31),
(160, 31),
(161, 31),
(163, 31),
(166, 31),
(167, 31),
(168, 31),
(174, 31),
(176, 31),
(179, 31),
(183, 31),
(186, 31),
(191, 31),
(194, 31),
(196, 31),
(197, 31),
(198, 31),
(101, 32),
(102, 32),
(112, 32),
(114, 32),
(120, 32),
(121, 32),
(124, 32),
(125, 32),
(127, 32),
(128, 32),
(129, 32),
(131, 32),
(132, 32),
(133, 32),
(135, 32),
(139, 32),
(143, 32),
(144, 32),
(150, 32),
(152, 32),
(154, 32),
(155, 32),
(156, 32),
(158, 32),
(159, 32),
(162, 32),
(165, 32),
(169, 32),
(170, 32),
(172, 32),
(173, 32),
(175, 32),
(177, 32),
(180, 32),
(181, 32),
(182, 32),
(184, 32),
(188, 32),
(189, 32),
(190, 32),
(193, 32),
(195, 32),
(199, 32),
(110, 33),
(101, 40),
(109, 40),
(112, 40),
(113, 40),
(115, 40),
(132, 40),
(135, 40),
(136, 40),
(140, 40),
(141, 40),
(147, 40),
(154, 40),
(177, 40),
(106, 41),
(108, 41),
(111, 41),
(122, 41),
(126, 41),
(131, 41),
(133, 41),
(137, 41),
(158, 41),
(160, 41),
(178, 41),
(180, 41),
(185, 41),
(188, 41),
(200, 41),
(118, 42),
(124, 42),
(145, 42),
(146, 42),
(162, 42),
(191, 42),
(192, 42),
(198, 42);

-- --------------------------------------------------------

--
-- Table structure for table `tags`
--

CREATE TABLE `tags` (
  `tag_id` bigint UNSIGNED NOT NULL,
  `name` varchar(100) NOT NULL,
  `type` varchar(50) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `tags`
--

INSERT INTO `tags` (`tag_id`, `name`, `type`, `created_at`, `updated_at`) VALUES
(1, 'Ayam', 'ingredient', NULL, NULL),
(2, 'Sapi', 'ingredient', NULL, NULL),
(3, 'Ikan', 'ingredient', NULL, NULL),
(4, 'Sayuran', 'ingredient', NULL, NULL),
(5, 'Nasi', 'ingredient', NULL, NULL),
(6, 'Mie', 'ingredient', NULL, NULL),
(7, 'Tahu', 'ingredient', NULL, NULL),
(8, 'Tempe', 'ingredient', NULL, NULL),
(10, 'Goreng', 'method', NULL, NULL),
(11, 'Bakar', 'method', NULL, NULL),
(12, 'Rebus', 'method', NULL, NULL),
(13, 'Tumis', 'method', NULL, NULL),
(14, 'Kuah', 'method', NULL, NULL),
(20, 'Pedas', 'taste', NULL, NULL),
(21, 'Manis', 'taste', NULL, NULL),
(22, 'Gurih', 'taste', NULL, NULL),
(23, 'Asam', 'taste', NULL, NULL),
(30, 'Sarapan', 'category', NULL, NULL),
(31, 'Makan Siang', 'category', NULL, NULL),
(32, 'Makan Malam', 'category', NULL, NULL),
(33, 'Dessert', 'category', NULL, NULL),
(40, 'Masakan Indonesia', 'cuisine', NULL, NULL),
(41, 'Masakan Jawa', 'cuisine', NULL, NULL),
(42, 'Masakan Sunda', 'cuisine', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `user_id` bigint UNSIGNED NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` enum('admin','editor','member') DEFAULT 'member',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`user_id`, `name`, `email`, `password`, `role`, `created_at`, `updated_at`) VALUES
(1, 'Admin Resep', 'admin@resep.com', '$2a$10$egUzvRpnZhJszx1BWjRW4eFpof8uBTJSnXCWWE3Rn338BhEuSfvFW', 'admin', '2025-10-17 19:32:29.621', '2025-10-17 19:32:29.621'),
(10, 'Budi Santoso', 'budi@contoh.com', '$2a$10$abcdefghijklmnopqrstuv', 'member', NULL, NULL),
(11, 'Ani Lestari', 'ani@contoh.com', '$2a$10$abcdefghijklmnopqrstuv', 'member', NULL, NULL),
(12, 'Candra Wijaya', 'candra@contoh.com', '$2a$10$abcdefghijklmnopqrstuv', 'member', NULL, NULL),
(13, 'Budi Baru', 'budibaru@contoh.com', '$2a$10$cO2/KRiZIIiGMj36tlJ04.FLGYKDJ0Za6nGwaZf3FaGsd/oTMgiPK', 'member', '2025-10-18 16:32:37.942', '2025-10-18 16:32:37.942');

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
(11, 101),
(12, 102),
(11, 111),
(12, 111),
(11, 113),
(11, 115),
(11, 116),
(12, 120),
(10, 121),
(12, 121),
(11, 126),
(10, 129),
(10, 135),
(10, 136),
(11, 136),
(12, 136),
(12, 143),
(10, 144),
(12, 144),
(11, 147),
(11, 148),
(12, 150),
(10, 151),
(12, 151),
(11, 157),
(10, 158),
(12, 159),
(10, 165),
(10, 166),
(11, 167),
(10, 173),
(12, 173),
(12, 174),
(11, 178),
(11, 179),
(11, 180),
(10, 181),
(11, 182),
(11, 188),
(12, 188),
(12, 189),
(11, 191),
(11, 192),
(10, 195),
(10, 196),
(11, 197),
(11, 198);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `menus`
--
ALTER TABLE `menus`
  ADD PRIMARY KEY (`menu_id`),
  ADD KEY `fk_menus_user` (`user_id`);

--
-- Indexes for table `menu_ratings`
--
ALTER TABLE `menu_ratings`
  ADD PRIMARY KEY (`rating_id`),
  ADD UNIQUE KEY `uni_menu_ratings_user_menu` (`user_id`,`menu_id`),
  ADD KEY `fk_menu_ratings_menu` (`menu_id`);

--
-- Indexes for table `menu_tags`
--
ALTER TABLE `menu_tags`
  ADD PRIMARY KEY (`menu_id`,`tag_id`),
  ADD KEY `fk_menu_tags_tag` (`tag_id`);

--
-- Indexes for table `tags`
--
ALTER TABLE `tags`
  ADD PRIMARY KEY (`tag_id`),
  ADD UNIQUE KEY `uni_tags_name` (`name`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `uni_users_email` (`email`);

--
-- Indexes for table `user_bookmarks`
--
ALTER TABLE `user_bookmarks`
  ADD PRIMARY KEY (`user_id`,`menu_id`),
  ADD KEY `fk_user_bookmarks_menu` (`menu_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `menus`
--
ALTER TABLE `menus`
  MODIFY `menu_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=201;

--
-- AUTO_INCREMENT for table `menu_ratings`
--
ALTER TABLE `menu_ratings`
  MODIFY `rating_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=65;

--
-- AUTO_INCREMENT for table `tags`
--
ALTER TABLE `tags`
  MODIFY `tag_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=43;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `user_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=14;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `menus`
--
ALTER TABLE `menus`
  ADD CONSTRAINT `fk_menus_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
  ADD CONSTRAINT `fk_users_menus` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

--
-- Constraints for table `menu_ratings`
--
ALTER TABLE `menu_ratings`
  ADD CONSTRAINT `fk_menu_ratings_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`),
  ADD CONSTRAINT `fk_menu_ratings_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

--
-- Constraints for table `menu_tags`
--
ALTER TABLE `menu_tags`
  ADD CONSTRAINT `fk_menu_tags_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`),
  ADD CONSTRAINT `fk_menu_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`tag_id`);

--
-- Constraints for table `user_bookmarks`
--
ALTER TABLE `user_bookmarks`
  ADD CONSTRAINT `fk_user_bookmarks_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`),
  ADD CONSTRAINT `fk_user_bookmarks_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
