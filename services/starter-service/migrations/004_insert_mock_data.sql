-- =============================================
-- INSERT MORE DEPARTMENTS
-- =============================================

INSERT INTO departments (group_department_id, business_unit_id, full_name, shortname, leader_id)
VALUES
-- VNGGames departments (business_unit_id = 1)
(5, 1, 'Game Analytics Team', 'GAT', NULL),
(5, 1, 'Game Security Team', 'GST', NULL),
(5, 1, 'Live Operations Team', 'LOT', NULL),
(5, 1, 'Player Support Team', 'PST', NULL),
(6, 1, 'Studio Alpha', 'SA', NULL),
(6, 1, 'Studio Beta', 'SB', NULL),
(6, 1, 'Studio Gamma', 'SG', NULL),
(7, 1, 'Mobile Games Development', 'MGD', NULL),
(7, 1, 'PC Games Development', 'PGD', NULL),
(7, 1, 'Cross-Platform Development', 'CPD', NULL),

-- Zalo departments (business_unit_id = 2)
(NULL, 2, 'Zalo Messenger Engineering', 'ZME', NULL),
(NULL, 2, 'Zalo Social Platform', 'ZSP', NULL),
(NULL, 2, 'Zalo AI Research', 'ZAI', NULL),
(NULL, 2, 'Zalo Infrastructure', 'ZIN', NULL),
(NULL, 2, 'Zalo Product Design', 'ZPD', NULL),
(NULL, 2, 'Zalo Quality Assurance', 'ZQA', NULL),

-- Zalopay departments (business_unit_id = 3)
(NULL, 3, 'Payment Gateway Team', 'PGT', NULL),
(NULL, 3, 'Risk Management Team', 'RMT', NULL),
(NULL, 3, 'Fraud Detection Team', 'FDT', NULL),
(NULL, 3, 'Financial Services Team', 'FST', NULL),
(NULL, 3, 'Merchant Solutions Team', 'MST', NULL),
(NULL, 3, 'Customer Experience Team', 'CXT', NULL),

-- Digital Business departments (business_unit_id = 4)
(NULL, 4, 'E-Commerce Platform', 'ECP', NULL),
(NULL, 4, 'Digital Advertising', 'DAD', NULL),
(NULL, 4, 'Cloud Services', 'CLS', NULL),
(NULL, 4, 'Enterprise Solutions', 'ENS', NULL),
(NULL, 4, 'Business Intelligence', 'BIT', NULL),
(NULL, 4, 'Digital Marketing', 'DMK', NULL);

-- =============================================
-- INSERT MORE STARTERS
-- =============================================

INSERT INTO starters (domain, name, email, mobile, work_phone, job_title, department_id, line_manager_id)
VALUES
-- Games Publishing Platform Engineering (dept 9)
('anhnt12', 'Nguyen Tuan Anh', 'anhnt12@vng.com.vn', '(+84) 0901234567', '0901234567', 'Senior Backend Engineer', 9, 8),
('binhpv', 'Pham Van Binh', 'binhpv@vng.com.vn', '(+84) 0902345678', '0902345678', 'Backend Engineer', 9, 8),
('chauld', 'Le Duc Chau', 'chauld@vng.com.vn', '(+84) 0903456789', '0903456789', 'DevOps Engineer', 9, 8),
('dungnt5', 'Nguyen Thanh Dung', 'dungnt5@vng.com.vn', '(+84) 0904567890', '0904567890', 'Frontend Engineer', 9, 8),
('emvt', 'Vo Thanh Em', 'emvt@vng.com.vn', '(+84) 0905678901', '0905678901', 'QA Engineer', 9, 8),

-- Game Data Studio (dept 12)
('gianght', 'Hoang Thu Giang', 'gianght@vng.com.vn', '(+84) 0906789012', '0906789012', 'Data Analyst', 12, 11),
('hieudv2', 'Dang Van Hieu', 'hieudv2@vng.com.vn', '(+84) 0907890123', '0907890123', 'Senior Data Scientist', 12, 11),
('oanhtpt', 'Phan Thi Oanh', 'oanhtpt@vng.com.vn', '(+84) 0908901234', '0908901234', 'Data Engineer', 12, 11),
('phucly', 'Ly Van Phuc', 'phucly@vng.com.vn', '(+84) 0909012345', '0909012345', 'ML Engineer', 12, 11),
('quynhct', 'Cao Thi Quynh', 'quynhct@vng.com.vn', '(+84) 0910123456', '0910123456', 'BI Developer', 12, 11),

-- Product Core (dept 8)
('sontv', 'Trinh Van Son', 'sontv@vng.com.vn', '(+84) 0911234567', '0911234567', 'Product Manager', 8, 6),
('thaodth', 'Duong Thi Thao', 'thaodth@vng.com.vn', '(+84) 0912345678', '0912345678', 'UX Designer', 8, 6),
('tuanhv', 'Ha Van Tuan', 'tuanhv@vng.com.vn', '(+84) 0913456789', '0913456789', 'Product Owner', 8, 6),
('uyenvt', 'Vu Thi Uyen', 'uyenvt@vng.com.vn', '(+84) 0914567890', '0914567890', 'UI Designer', 8, 6),
('vinhlv', 'Lam Van Vinh', 'vinhlv@vng.com.vn', '(+84) 0915678901', '0915678901', 'Business Analyst', 8, 6),

-- Platform Integration (dept 10)
('xuanct', 'Chau Thi Xuan', 'xuanct@vng.com.vn', '(+84) 0916789012', '0916789012', 'Integration Engineer', 10, 6),
('yennth', 'Nguyen Thi Yen', 'yennth@vng.com.vn', '(+84) 0917890123', '0917890123', 'API Developer', 10, 6),
('anhpv3', 'Pham Van Anh', 'anhpv3@vng.com.vn', '(+84) 0918901234', '0918901234', 'System Integrator', 10, 6),
('baotn', 'Tran Ngoc Bao', 'baotn@vng.com.vn', '(+84) 0919012345', '0919012345', 'Technical Lead', 10, 6),
('cuongdv', 'Dang Van Cuong', 'cuongdv@vng.com.vn', '(+84) 0920123456', '0920123456', 'Solutions Architect', 10, 6),

-- Game Infrastructure Operation (dept 11)
('datnq', 'Nguyen Quoc Dat', 'datnq@vng.com.vn', '(+84) 0921234567', '0921234567', 'DevOps Lead', 11, 7),
('emhtt', 'Hoang Thi Em', 'emhtt@vng.com.vn', '(+84) 0922345678', '0922345678', 'Site Reliability Engineer', 11, 7),
('fongbt', 'Bui Thanh Fong', 'fongbt@vng.com.vn', '(+84) 0923456789', '0923456789', 'Cloud Engineer', 11, 7),
('giangnh', 'Nguyen Hoang Giang', 'giangnh@vng.com.vn', '(+84) 0924567890', '0924567890', 'Infrastructure Engineer', 11,
 7),
('haivl', 'Vu Loc Hai', 'haivl@vng.com.vn', '(+84) 0925678901', '0925678901', 'Network Engineer', 11, 7),

-- Zalo Messenger Engineering
('inhpd', 'Pham Duc Inh', 'inhpd@vng.com.vn', '(+84) 0926789012', '0926789012', 'Backend Engineer', NULL, 5),
('khantn', 'Tran Ngoc Khan', 'khantn@vng.com.vn', '(+84) 0927890123', '0927890123', 'Mobile Developer', NULL, 5),
('lanbt', 'Bui Thanh Lan', 'lanbt@vng.com.vn', '(+84) 0928901234', '0928901234', 'iOS Developer', NULL, 5),
('minhnq2', 'Nguyen Quoc Minh', 'minhnq2@vng.com.vn', '(+84) 0929012345', '0929012345', 'Android Developer', NULL, 5),
('nampv2', 'Pham Van Nam', 'nampv2@vng.com.vn', '(+84) 0930123456', '0930123456', 'Full Stack Engineer', NULL, 5),

-- Zalopay Payment Gateway
('oanhht', 'Hoang Thi Oanh', 'oanhht@vng.com.vn', '(+84) 0931234567', '0931234567', 'Payment Engineer', NULL, 3),
('phongnd', 'Nguyen Duc Phong', 'phongnd@vng.com.vn', '(+84) 0932345678', '0932345678', 'Security Engineer', NULL, 3),
('quanpv', 'Pham Van Quan', 'quanpv@vng.com.vn', '(+84) 0933456789', '0933456789', 'Backend Developer', NULL, 3),
('rubylt', 'Le Thi Ruby', 'rubylt@vng.com.vn', '(+84) 0934567890', '0934567890', 'Risk Analyst', NULL, 3),
('sonnh2', 'Nguyen Hoang Son', 'sonnh2@vng.com.vn', '(+84) 0935678901', '0935678901', 'Fraud Detection Specialist',
 NULL, 3),

-- Digital Business
('thanhnt3', 'Nguyen Thanh Thanh', 'thanhnt3@vng.com.vn', '(+84) 0936789012', '0936789012', 'Solutions Engineer', NULL,
 4),
('uyenpd', 'Pham Duc Uyen', 'uyenpd@vng.com.vn', '(+84) 0937890123', '0937890123', 'Cloud Architect', NULL, 4),
('vinhlt', 'Le Thanh Vinh', 'vinhlt@vng.com.vn', '(+84) 0938901234', '0938901234', 'Enterprise Consultant', NULL, 4),
('xuannh', 'Nguyen Hoang Xuan', 'xuannh@vng.com.vn', '(+84) 0939012345', '0939012345', 'Technical Sales', NULL, 4),
('yenpth', 'Phan Thi Yen', 'yenpth@vng.com.vn', '(+84) 0940123456', '0940123456', 'Account Manager', NULL, 4),

-- More Game Studio employees
('anhtv2', 'Tran Van Anh', 'anhtv2@vng.com.vn', '(+84) 0941234567', '0941234567', 'Game Designer', NULL, 2),
('binhnv', 'Nguyen Van Binh', 'binhnv@vng.com.vn', '(+84) 0942345678', '0942345678', 'Level Designer', NULL, 2),
('chiptt', 'Pham Thi Chi', 'chiptt@vng.com.vn', '(+84) 0943456789', '0943456789', 'Game Artist', NULL, 2),
('ducnv2', 'Nguyen Van Duc', 'ducnv2@vng.com.vn', '(+84) 0944567890', '0944567890', '3D Artist', NULL, 2),
('emlt', 'Le Thi Em', 'emlt@vng.com.vn', '(+84) 0945678901', '0945678901', 'Animator', NULL, 2),

-- Executive Services team
('fongpv', 'Pham Van Fong', 'fongpv@vng.com.vn', '(+84) 0946789012', '0946789012', 'Executive Assistant', 3, 1),
('gianglt', 'Le Thi Giang', 'gianglt@vng.com.vn', '(+84) 0947890123', '0947890123', 'HR Manager', 3, 1),
('hanh', 'Nguyen Thi Hanh', 'hanh@vng.com.vn', '(+84) 0948901234', '0948901234', 'Admin Manager', 3, 1),
('inhnt', 'Nguyen Thanh Inh', 'inhnt@vng.com.vn', '(+84) 0949012345', '0949012345', 'Finance Manager', 3, 1),
('kimbui', 'Bui Thi Kim', 'kimbui@vng.com.vn', '(+84) 0950123456', '0950123456', 'Legal Counsel', 3, 1);

INSERT INTO departments (group_department_id, business_unit_id, full_name, shortname, leader_id)
VALUES
-- Under Games Publishing Platform Engineering (dept 5) - Already has some children
(5, 1, 'Game Monetization Team', 'GMT', NULL),
(5, 1, 'Publishing Operations', 'POP', NULL),
(5, 1, 'Community Management', 'CMT', NULL),

-- Under Game Studio (dept 6) - Already has Studio Alpha, Beta, Gamma
(6, 1, 'Creative Services', 'CRS', NULL),
(6, 1, 'QA & Testing', 'QAT', NULL),

-- Under Platform Engineering (dept 7) - Already has Mobile, PC, Cross-Platform
(7, 1, 'Game Engine Team', 'GET', NULL),
(7, 1, 'Tools & SDK', 'TSD', NULL),

-- Under Product Core (dept 8)
(8, 1, 'Product Strategy', 'PST', NULL),
(8, 1, 'User Experience', 'UXP', NULL),
(8, 1, 'Product Analytics', 'PAN', NULL),

-- Under Games Publishing Platform Engineering > Games Platform Backend (dept 9)
(9, 1, 'Backend Infrastructure', 'BIN', NULL),
(9, 1, 'API Gateway', 'AGW', NULL),
(9, 1, 'Microservices', 'MSV', NULL),

-- Under Product Core > Platform Integration (dept 10)
(10, 1, 'Third-Party Integration', 'TPI', NULL),
(10, 1, 'Payment Integration', 'PIN', NULL),

-- Under Platform Engineering > Game Infrastructure Operation (dept 11)
(11, 1, 'Cloud Infrastructure', 'CIN', NULL),
(11, 1, 'Network Operations', 'NOP', NULL),
(11, 1, 'Monitoring & Alerting', 'MNA', NULL),

-- Under Platform Engineering > Game Data Studio (dept 12)
(12, 1, 'Data Pipeline', 'DPL', NULL),
(12, 1, 'Analytics Platform', 'APL', NULL),
(12, 1, 'Machine Learning', 'MLT', NULL);


INSERT INTO departments (group_department_id, business_unit_id, full_name, shortname, leader_id)
VALUES
-- Under Studio Alpha (dept 17)
(17, 1, 'Alpha - Game Design', 'AGD', NULL),
(17, 1, 'Alpha - Programming', 'APG', NULL),
(17, 1, 'Alpha - Art & Animation', 'AAA', NULL),

-- Under Studio Beta (dept 18)
(18, 1, 'Beta - Game Design', 'BGD', NULL),
(18, 1, 'Beta - Programming', 'BPG', NULL),

-- Under Mobile Games Development (dept 20)
(20, 1, 'iOS Development', 'IOS', NULL),
(20, 1, 'Android Development', 'AND', NULL),
(20, 1, 'Mobile QA', 'MQA', NULL),

-- Under PC Games Development (dept 21)
(21, 1, 'Windows Client', 'WIN', NULL),
(21, 1, 'Game Optimization', 'GOP', NULL),

-- Under Backend Infrastructure (will be dept 31)
(31, 1, 'Database Team', 'DBT', NULL),
(31, 1, 'Cache & Storage', 'CAS', NULL),

-- Under Cloud Infrastructure (will be dept 39)
(39, 1, 'AWS Operations', 'AWS', NULL),
(39, 1, 'Container Platform', 'CPF', NULL),

-- Under Data Pipeline (will be dept 42)
(42, 1, 'ETL Development', 'ETL', NULL),
(42, 1, 'Data Quality', 'DQA', NULL);

INSERT INTO departments (group_department_id, business_unit_id, full_name, shortname, leader_id)
VALUES
-- Under Alpha - Programming (dept 49)
(49, 1, 'Gameplay Programming', 'GPP', NULL),
(49, 1, 'Engine Programming', 'EPG', NULL),
(49, 1, 'Network Programming', 'NPG', NULL),

-- Under iOS Development (dept 53)
(53, 1, 'iOS UI/UX', 'IUX', NULL),
(53, 1, 'iOS Backend Integration', 'IBI', NULL),

-- Under Android Development (dept 54)
(54, 1, 'Android UI/UX', 'AUX', NULL),
(54, 1, 'Android Performance', 'APF', NULL),

-- Under Database Team (dept 58)
(58, 1, 'MySQL Administration', 'MYA', NULL),
(58, 1, 'MongoDB Operations', 'MGO', NULL);

-- VNG Corporation (id=1)
-- ├─ VNGGames (id=2)
-- │  ├─ Executive Services (id=3)
-- │  ├─ Games Publishing Platform Engineering (id=5)
-- │  │  ├─ Game Analytics Team (id=13)
-- │  │  ├─ Game Security Team (id=14)
-- │  │  ├─ Live Operations Team (id=15)
-- │  │  ├─ Player Support Team (id=16)
-- │  │  ├─ Game Monetization Team (id=25)
-- │  │  ├─ Publishing Operations (id=26)
-- │  │  ├─ Community Management (id=27)
-- │  │  └─ Games Platform Backend (id=9)
-- │  │     ├─ Backend Infrastructure (id=31)
-- │  │     │  ├─ Database Team (id=58)
-- │  │     │  │  ├─ MySQL Administration (id=65)
-- │  │     │  │  └─ MongoDB Operations (id=66)
-- │  │     │  └─ Cache & Storage (id=59)
-- │  │     ├─ API Gateway (id=32)
-- │  │     └─ Microservices (id=33)
-- │  ├─ Game Studio (id=6)
-- │  │  ├─ Studio Alpha (id=17)
-- │  │  │  ├─ Alpha - Game Design (id=47)
-- │  │  │  ├─ Alpha - Programming (id=48)
-- │  │  │  │  ├─ Gameplay Programming (id=61)
-- │  │  │  │  ├─ Engine Programming (id=62)
-- │  │  │  │  └─ Network Programming (id=63)
-- │  │  │  └─ Alpha - Art & Animation (id=49)
-- │  │  ├─ Studio Beta (id=18)
-- │  │  │  ├─ Beta - Game Design (id=50)
-- │  │  │  └─ Beta - Programming (id=51)
-- │  │  ├─ Studio Gamma (id=19)
-- │  │  ├─ Creative Services (id=28)
-- │  │  └─ QA & Testing (id=29)
-- │  ├─ Platform Engineering (id=7)
-- │  │  ├─ Mobile Games Development (id=20)
-- │  │  │  ├─ iOS Development (id=52)
-- │  │  │  │  ├─ iOS UI/UX (id=67)
-- │  │  │  │  └─ iOS Backend Integration (id=68)
-- │  │  │  ├─ Android Development (id=53)
-- │  │  │  │  ├─ Android UI/UX (id=69)
-- │  │  │  │  └─ Android Performance (id=70)
-- │  │  │  └─ Mobile QA (id=54)
-- │  │  ├─ PC Games Development (id=21)
-- │  │  │  ├─ Windows Client (id=55)
-- │  │  │  └─ Game Optimization (id=56)
-- │  │  ├─ Cross-Platform Development (id=22)
-- │  │  ├─ Game Engine Team (id=30)
-- │  │  ├─ Tools & SDK (id=31)
-- │  │  ├─ Game Infrastructure Operation (id=11)
-- │  │  │  ├─ Cloud Infrastructure (id=39)
-- │  │  │  │  ├─ AWS Operations (id=60)
-- │  │  │  │  └─ Container Platform (id=61)
-- │  │  │  ├─ Network Operations (id=40)
-- │  │  │  └─ Monitoring & Alerting (id=41)
-- │  │  └─ Game Data Studio (id=12)
-- │  │     ├─ Data Pipeline (id=42)
-- │  │     │  ├─ ETL Development (id=71)
-- │  │     │  └─ Data Quality (id=72)
-- │  │     ├─ Analytics Platform (id=43)
-- │  │     └─ Machine Learning (id=44)
-- │  └─ Product Core (id=8)
-- │     ├─ Product Strategy (id=35)
-- │     ├─ User Experience (id=36)
-- │     ├─ Product Analytics (id=37)
-- │     └─ Platform Integration (id=10)
-- │        ├─ Third-Party Integration (id=38)
-- │        └─ Payment Integration (id=39)